package load

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

func NewMinioClient() *minio.Client {
	endpoint := "localhost:9000"
	accessKeyID := "pQVJVMIAidUiyYH2"
	secretAccessKey := "dtTpbfJZpjuUT6KcTHWwekhUT0FxI4Xx"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""), Secure: useSSL})

	if err != nil {
		logrus.Fatalln(err)
	}

	return minioClient
}
