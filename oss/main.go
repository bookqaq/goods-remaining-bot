package oss

import (
	"context"

	"github.com/minio/minio-go/v7"
)

const (
	Bucket_name = "goods-remaining-bot-image-bucket"
)

func UploadFile(name string, path string, mtype string) error {
	_, err := Client.FPutObject(context.Background(), Bucket_name,
		name, path, minio.PutObjectOptions{ContentType: mtype})
	return err
}
