package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

const (
	secret        = "password"
	key           = "user"
	repopw        = "password"
	repo          = "s3:http://localhost:9000/test"
	resticVersion = "v0.9.5"
	resticRepo    = "https://github.com/restic/restic/releases/download"
)

func TestIntegration(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "wrestic")

	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(tmpdir)
	killFunc := prepare(t, tmpdir)
	runIntegrationTests(t, tmpdir)
	teardown(t, killFunc)
}

func prepare(t *testing.T, tmpdir string) func() {
	fmt.Println("=================== Preparing test env ===================")

	err := os.Mkdir(filepath.Join(tmpdir, "bin"), 0777)

	if err != nil {
		panic(err)
	}

	getMinio(tmpdir)
	getRestic(tmpdir)
	return runMinio(t, tmpdir)
}

func teardown(t *testing.T, f func()) {
	fmt.Println("=================== Tearing down test env ===================")
	f()
}

func runIntegrationTests(t *testing.T, tmpdir string) {
	fmt.Println("=================== Starting tests ===================")
	cmd := exec.Command("go", "test", "-v", "-mod", "vendor", "-tags", "integration", "./cmd/wrestic/...")
	resticBin, _ := filepath.Abs(filepath.Join(tmpdir, "bin", "restic"))
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
	addOutputs(cmd)
	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("=================== Tests finished ===================")
}

func getMinio(tmpdir string) {

	fmt.Println("Downloading minio...")

	minioPath := filepath.Join(tmpdir, "bin", "minio")

	if _, err := os.Stat(minioPath); os.IsExist(err) {
		fmt.Println("Minio already downloaded")
		return
	}

	url := fmt.Sprintf("https://dl.minio.io/server/minio/release/%s-%s/minio",
		runtime.GOOS, runtime.GOARCH)

	downloadBinary(minioPath, url)

}

func runMinio(t *testing.T, tmpdir string) func() {
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

func getRestic(tmpdir string) {
	fmt.Println("Downloading restic")

	bzipPath := filepath.Join(tmpdir, "bin", "restic.bz2")

	cleanVer := strings.Replace(resticVersion, "v", "", 1)
	url := fmt.Sprintf("%s/%s/restic_%s_%s_%s.bz2",
		resticRepo, resticVersion, cleanVer, runtime.GOOS, runtime.GOARCH)

	err := downloadBinary(bzipPath, url)

	if err != nil {
		panic(err)
	}

	extract := exec.Command("bzip2", "-d", bzipPath)
	addOutputs(extract)
	extract.Dir = filepath.Dir(bzipPath)
	extract.Run()

}

func addOutputs(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func downloadBinary(binPath string, url string) error {
	if _, err := os.Stat(binPath); os.IsExist(err) {
		fmt.Println("Binary already downloaded")
		return nil
	}

	binFile, err := os.Create(binPath)
	if err != nil {
		return fmt.Errorf("create tempfile for download failed: %v", err)
	}
	fmt.Printf("downloading %v\n", url)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading: %v", err)
	}

	_, err = io.Copy(binFile, res.Body)
	if err != nil {
		return fmt.Errorf("error saving to file: %v", err)
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing HTTP download: %v", err)
	}

	err = binFile.Close()
	if err != nil {
		fmt.Printf("closing tempfile failed: %v\n", err)
		return fmt.Errorf("error closing file: %v", err)
	}

	err = os.Chmod(binFile.Name(), 0755)
	if err != nil {
		return fmt.Errorf("chmod() failed: %v", err)
	}

	fmt.Printf("downloaded to %v\n", binFile.Name())
	return nil
}
