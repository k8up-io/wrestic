package rest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
func PodExec(params Params) (io.Reader, error) {
	config, _ := getClientConfig()
	k8sclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req := k8sclient.Core().RESTClient().Post().
		Resource("pods").
		Name(params.Pod).
		Namespace(params.Namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := apiv1.AddToScheme(scheme); err != nil {
		fmt.Println(err)
		return nil, err
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&apiv1.PodExecOptions{
		Command:   strings.Fields(params.BackupCommand),
		Container: params.Container,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		fmt.Println(err)
		return nil, err
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
			// return nil, err
		}
		stdoutWriter.Close()
	}()

	return stdoutReader, nil
}

func getClientConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		err1 := err
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			err = fmt.Errorf("InClusterConfig as well as BuildConfigFromFlags Failed. Error in InClusterConfig: %+v\nError in BuildConfigFromFlags: %+v", err1, err)
			return nil, err
		}
	}

	return config, nil
}
