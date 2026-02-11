package s3

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type Uploader struct {
	client     *minio.Client
	bucket     string
	publicBase string
}

func NewUploader(client *minio.Client, bucket, publicBase string) *Uploader {
	return &Uploader{client: client, bucket: bucket, publicBase: strings.TrimRight(publicBase, "/")}
}

func (u *Uploader) UploadMany(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	urls := make([]string, 0, len(files))

	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			return nil, err
		}

		ext := strings.ToLower(path.Ext(fh.Filename))
		if ext == "" {
			ext = ".bin"
		}

		key := fmt.Sprintf("blogs/%s%s", randHex(16), ext)

		_, err = u.client.PutObject(ctx, u.bucket, key, f, fh.Size, minio.PutObjectOptions{
			ContentType: fh.Header.Get("Content-Type"),
		})
		f.Close()
		if err != nil {
			return nil, err
		}

		urls = append(urls, fmt.Sprintf("%s/%s", u.publicBase, key))
	}

	return urls, nil
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b) + fmt.Sprintf("-%d", time.Now().UnixNano())
}
