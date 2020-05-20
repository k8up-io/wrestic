package kubernetes

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/firepear/qsplit"
	"github.com/go-logr/logr"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
)

type ExecData struct {
	Reader *io.PipeReader
	Done   chan bool
}

// PodExec sends the command to the specified pod
// and returns a bytes buffer with the stdout
func PodExec(pod BackupPod, log logr.Logger) (*ExecData, *bytes.Buffer, error) {

	execLogger := log.WithName("k8sExec")

	config, _ := getClientConfig()
	k8sclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create k8s for exec: %v", err)
	}

	req := k8sclient.Core().RESTClient().Post().
		Resource("pods").
		Name(pod.PodName).
		Namespace(pod.Namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := apiv1.AddToScheme(scheme); err != nil {
		return nil, nil, fmt.Errorf("can't add runtime scheme: %v", err)
	}

	command := qsplit.ToStrings([]byte(pod.Command))
	execLogger.Info("executing command", "command", strings.Join(command, ", "), "namespace", pod.Namespace, "pod", pod.PodName)

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&apiv1.PodExecOptions{
		Command:   command,
		Container: pod.ContainerName,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, nil, err
	}

	var stderr bytes.Buffer
	var stdoutReader, stdoutWriter = io.Pipe()
	done := make(chan bool, 1)
	go func() {
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: stdoutWriter,
			Stderr: &stderr,
			Tty:    false,
		})

		defer stdoutWriter.Close()
		done <- true

		if err != nil {
			execLogger.Error(err, "error ocurred stream backup data", "namespace", pod.Namespace, "pod", pod.PodName)
			return
		}
	}()

	data := &ExecData{
		Done:   done,
		Reader: stdoutReader,
	}

	return data, &stderr, nil
}
