package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

const (
	tmpdir = "tmp"
	secret = "password"
	key    = "user"
	repopw = "password"
	repo   = "s3:http://localhost:9000/test"
)

func TestIntegration(t *testing.T) {
	killFunc := prepare(t)
	runIntegrationTests(t)
	teardown(t, killFunc)
}

func prepare(t *testing.T) func() {
	fmt.Println("=================== Preparing test env ===================")
	if _, ok := os.Stat(tmpdir); os.IsExist(ok) {
		os.RemoveAll(tmpdir)
	}
	getMinio()
	buildCustomRestic(t)
	return runMinio(t)
}

func teardown(t *testing.T, f func()) {
	fmt.Println("=================== Tearing down test env ===================")
	f()
	os.RemoveAll(tmpdir)
}

func runIntegrationTests(t *testing.T) {
	fmt.Println("=================== Starting tests ===================")
	cmd := exec.Command("go", "test", "-v", "-tags", "integration", "./cmd/wrestic/...")
	resticBin, _ := filepath.Abs(filepath.Join(tmpdir, "restic/restic"))
	fmt.Println("Restic location", resticBin)
	cmd.Env = append(os.Environ(),
		"RESTIC_PASSWORD="+repopw,
		"RESTIC_REPOSITORY="+repo,
		"AWS_SECRET_ACCESS_KEY="+secret,
		"AWS_ACCESS_KEY_ID="+key,
		"BACKUP_DIR=testdata",
		"STATS_URL=http://localhost:8091",
		"RESTORE_ACCESSKEYID="+key,
		"RESTORE_SECRETACCESSKEY="+secret,
		"RESTORE_S3ENDPOINT=http://localhost:9000/restore",
		"RESTIC_BINARY="+resticBin,
		"RESTORE_DIR="+filepath.Join(tmpdir, "restore"),
	)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		t.Error(err)
	}
	fmt.Println("=================== Tests finished ===================")
}

func getMinio() error {
	minioPath := filepath.Join(os.Getenv("GOPATH"), "bin", "minio")

	if _, err := os.Stat(minioPath); !os.IsNotExist(err) {
		fmt.Println("Minio already downloaded")
		return nil
	}

	tempfile, err := os.Create(minioPath)
	if err != nil {
		return fmt.Errorf("create tempfile for minio download failed: %v", err)
	}

	url := fmt.Sprintf("https://dl.minio.io/server/minio/release/%s-%s/archive/minio.RELEASE.2019-03-06T22-47-10Z",
		runtime.GOOS, runtime.GOARCH)
	fmt.Printf("downloading %v\n", url)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading minio server: %v", err)
	}

	_, err = io.Copy(tempfile, res.Body)
	if err != nil {
		return fmt.Errorf("error saving minio server to file: %v", err)
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing HTTP download: %v", err)
	}

	err = tempfile.Close()
	if err != nil {
		fmt.Printf("closing tempfile failed: %v\n", err)
		return fmt.Errorf("error closing minio server file: %v", err)
	}

	err = os.Chmod(tempfile.Name(), 0755)
	if err != nil {
		return fmt.Errorf("chmod(minio-server) failed: %v", err)
	}

	fmt.Printf("downloaded minio server to %v\n", tempfile.Name())
	return nil
}

func runMinio(t *testing.T) func() {
	configDir := filepath.Join(tmpdir, "config")
	rootDir := filepath.Join(tmpdir, "root")

	os.MkdirAll(configDir, 0700)
	os.MkdirAll(rootDir, 0700)

	cmd := exec.Command("minio",
		"server",
		"--address", "127.0.0.1:9000",
		"--config-dir", configDir,
		rootDir)
	cmd.Env = append(os.Environ(), getMinioEnv()...)
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	// wait until the TCP port is reachable
	var success bool
	for i := 0; i < 100; i++ {
		time.Sleep(200 * time.Millisecond)

		c, err := net.Dial("tcp", "localhost:9000")
		if err == nil {
			success = true
			if err := c.Close(); err != nil {
				t.Fatal(err)
			}
			break
		}
	}

	if !success {
		t.Fatal("unable to connect to minio server")
		return nil
	}

	return func() {
		err = cmd.Process.Kill()
		if err != nil {
			t.Fatal(err)
		}

		// ignore errors, we've killed the process
		_ = cmd.Wait()
	}
}

func getMinioEnv() []string {
	return []string{
		"MINIO_ACCESS_KEY=" + key,
		"MINIO_SECRET_KEY=" + secret,
	}
}

func buildCustomRestic(t *testing.T) {
	fmt.Println("Cloning restic")
	clone := exec.Command("git", "clone", "https://github.com/Kidswiss/restic", "tmp/restic")

	out, err := clone.CombinedOutput()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(out))

	checkout := exec.Command("git", "checkout", "tar")
	checkout.Dir = "tmp/restic"
	checkoutOut, err := checkout.CombinedOutput()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(checkoutOut))

	fmt.Println("Building restic")

	build := exec.Command("go", "run", "-mod=vendor", "build.go", "-v")
	build.Dir = "tmp/restic"

	buildOut, err := build.CombinedOutput()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(buildOut))
}
