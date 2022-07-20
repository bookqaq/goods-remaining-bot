package oss

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	Endpoint  = "your_minio_oss" // must NOT use https url
	accessKey = "your_key"
	secretKey = "your_secret"
	useSSL    = true
)

var Client *minio.Client

func Connect() {
	client, err := minio.New(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Panic(err)
	}

	Client = client
	log.Println("Minio connected.")
}
