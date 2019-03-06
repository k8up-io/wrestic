package s3

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	minio "github.com/minio/minio-go"
)

// Client wraps the minio s3 client
type Client struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	minioClient     *minio.Client
	bucket          string
}

type UploadObject struct {
	ObjectStream io.Reader
	Name         string
}

// New returns a new Client
func New(endpoint, accessKeyID, secretAccessKey string) *Client {
	return &Client{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
}

// Connect creates a minio client
func (c *Client) Connect() error {

	u, err := url.Parse(c.Endpoint)

	var ssl bool
	if u.Scheme == "https" {
		ssl = true
	} else if u.Scheme == "http" {
		ssl = false
	} else {
		return fmt.Errorf("Endpoint %v has wrong scheme (http/https)", c.Endpoint)
	}

	c.bucket = strings.Replace(u.Path, "/", "", 1)
	c.Endpoint = u.Host
	mc, err := minio.New(c.Endpoint, c.AccessKeyID, c.SecretAccessKey, ssl)
	c.minioClient = mc

	if err == nil {
		err = c.createBucket()
	}

	return err
}

func (c *Client) createBucket() error {
	exists, err := c.minioClient.BucketExists(c.bucket)
	// Workaround for upstream bug -> australian s3 returns error on non existing bucket.
	if !exists && (err == nil || strings.Contains(err.Error(), "exist")) {
		return c.minioClient.MakeBucket(c.bucket, "")
	} else if err != nil {
		return err
	}
	return nil
}

// Upload uploads a io.Reader object to the configured endpoint
func (c *Client) Upload(object UploadObject) error {
	_, err := c.minioClient.PutObject(c.bucket, object.Name, object.ObjectStream, -1, minio.PutObjectOptions{})
	return err
}

// Get gets a file or returns an error.
func (c *Client) Get(filename string) (*minio.Object, error) {
	return c.minioClient.GetObject(c.bucket, filename, minio.GetObjectOptions{})
}

// Stat returns metainformation about an object in the repository.
func (c *Client) Stat(filename string) (minio.ObjectInfo, error) {
	return c.minioClient.StatObject(c.bucket, filename, minio.StatObjectOptions{})
}
