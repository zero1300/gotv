package load

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

func NewMinioClient() *minio.Client {
	endpoint := "192.168.10.11:9000"
	accessKeyID := "jRDw9HsWCs6DbZIe"
	secretAccessKey := "61wDepkgP6TqI3dE0aCO2I4kpMtX4T0w"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""), Secure: useSSL})

	if err != nil {
		logrus.Fatalln(err)
	}

	return minioClient
}
