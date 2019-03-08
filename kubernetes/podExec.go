package kubernetes

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"firepear.net/qsplit"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
)

// Params contains the necessary parameters
// for the PodExec function
type Params struct {
	BackupCommand string
	Pod           string
	Container     string
	Namespace     string
}

// PodExec sends the command to the specified pod
// and returns a bytes buffer with the stdout
func PodExec(params Params) (io.Reader, *bytes.Buffer, error) {
	config, _ := getClientConfig()
	k8sclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create k8s for exec: %v", err)
	}

	req := k8sclient.Core().RESTClient().Post().
		Resource("pods").
		Name(params.Pod).
		Namespace(params.Namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := apiv1.AddToScheme(scheme); err != nil {
		return nil, nil, fmt.Errorf("can't add runtime scheme: %v", err)
	}

	command := qsplit.ToStrings([]byte(params.BackupCommand))
	fmt.Printf("Backup command: %v\n", strings.Join(command, ", "))

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&apiv1.PodExecOptions{
		Command:   command,
		Container: params.Container,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	var stderr bytes.Buffer
	var stdoutReader, stdoutWriter = io.Pipe()
	go func() {
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: stdoutWriter,
			Stderr: &stderr,
			Tty:    false,
		})

		if err != nil {
			fmt.Println(err)
			stdoutReader.CloseWithError(err)
			stdoutWriter.CloseWithError(err)
			return
		}
		stdoutWriter.Close()
	}()

	return stdoutReader, &stderr, nil
}
