package restic

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"git.vshn.net/vshn/wrestic/kubernetes"
	"git.vshn.net/vshn/wrestic/output"
	"github.com/prometheus/client_golang/prometheus"
)

type genericCommand struct {
	errorMessage      error
	stdOut, stdErrOut []string
	// command is the actual instance of the command that is handled by this
	// instance of genericCommand
	command *exec.Cmd
	// commandState holds the global state what command is currently running
	commandState *commandState
}

type commandOptions struct {
	print bool
	stdin bool
	kubernetes.Params
	// output can be used to get realtime output from the restic command
	// not all subcommands have realtime output though, so it's not mandatory
	// to set it.
	output chan string
}

func newGenericCommand(commandState *commandState) *genericCommand {
	return &genericCommand{
		stdOut:       make([]string, 0),
		stdErrOut:    make([]string, 0),
		commandState: commandState,
	}
}

func (g *genericCommand) exec(args []string, options commandOptions) {

	cmd := exec.Command(getResticBin(), args...)
	cmd.Env = os.Environ()

	g.command = cmd

	if options.stdin {
		stdout, stderr, err := kubernetes.PodExec(options.Params)
		if err != nil {
			fmt.Println(err)
			g.errorMessage = err
			return
		}
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
			g.errorMessage = err
			return
		}
		if stdout == nil {
			fmt.Println("stdout is nil")
		}
		// This needs to run in a separate thread because
		// cmd.CombinedOutput blocks until the command is finished
		// TODO: this is the place where we could implement some sort of
		// progress bars by wrapping stdin/stdout in a custom reader/writer
		go func() {
			defer stdin.Close()
			_, err := io.Copy(stdin, stdout)
			if err != nil {
				cmd.Process.Kill()
				fmt.Println(err)
				g.errorMessage = err
				stderrStr := stderr.String()
				if stderrStr != "" {
					fmt.Printf("Stderr of pod exec: '%v'", stderr)
					g.errorMessage = errors.New(stderrStr)
				}
			}
		}()
	}

	commandStdout, err := cmd.StdoutPipe()
	commandStderr, err := cmd.StderrPipe()

	finished := make(chan error, 0)

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		g.errorMessage = err
		return
	}

	g.commandState.running = g

	go func() {
		var collectErr error
		g.stdOut, collectErr = g.collectOutput(commandStdout, options.print, options.output)
		finished <- collectErr
	}()

	go func() {
		var collectErr error
		g.stdErrOut, collectErr = g.collectOutput(commandStderr, options.print, options.output)
		finished <- collectErr
	}()

	collectErr1 := <-finished
	collectErr2 := <-finished
	err = cmd.Wait()

	// Avoid overwriting any errors produced by the
	// copy command
	if g.errorMessage == nil {
		if err != nil {
			g.errorMessage = err
		}
		if collectErr1 != nil {
			g.errorMessage = collectErr1
		}
		if collectErr2 != nil {
			g.errorMessage = collectErr2
		}
	}
}

func (g *genericCommand) collectOutput(output io.Reader, print bool, out chan string) ([]string, error) {
	collectedOutput := make([]string, 0)
	scanner := bufio.NewScanner(output)
	buff := make([]byte, 64*1024*1024)
	scanner.Buffer(buff, 64*1024*1024)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		if print {
			fmt.Println(m)
		}
		collectedOutput = append(collectedOutput, m)
		if out != nil {
			out <- m
		}
	}
	return collectedOutput, scanner.Err()
}

// GetError returns if there was an error
func (g *genericCommand) GetError() error { return g.errorMessage }

// GetStdOut returns the complete output of the command
func (g *genericCommand) GetStdOut() []string { return g.stdOut }

// GetStdErrOut returns the complete StdErr of the command
func (g *genericCommand) GetStdErrOut() []string { return g.stdErrOut }

// GetWebhookData returns all objects that should get marshalled to json and
// sent to the webhook endpoint. Returns nil by default.
func (g *genericCommand) GetWebhookData() []output.JsonMarshaller {
	return nil
}

// ToProm returns a list of prometheus collectors that should get pushed to
// the prometheus push gateway.
func (g *genericCommand) ToProm() []prometheus.Collector {
	return nil
}

func (g *genericCommand) sendSignal(signal os.Signal) error {
	return g.command.Process.Signal(signal)
}
