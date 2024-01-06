package storage

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type StorageS3Stuct struct {
	StorageS3  *session.Session
	BucketName string
}

func NewS3() *StorageS3Stuct {

	key := os.Getenv("AWS_KEY")
	secret := os.Getenv("AWS_SECRET")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_BUCKET")
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			key,    // id
			secret, // secret
			""),    // token can be left blank for now
	})

	if err != nil {
		return nil
	}

	fmt.Println(s)

	return &StorageS3Stuct{
		StorageS3:  s,
		BucketName: bucket,
	}
}

func ConnectAws() *s3.S3 {

	aws_access_key_id := os.Getenv("AWS_ACCESS_KEY_ID")
	aws_secret_access_key := os.Getenv("AWS_SECRET")
	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
	_, err := creds.Get()
	if err != nil {
		// handle error
		fmt.Print(err)

	}
	cfg := aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")).WithCredentials(creds)
	svc := s3.New(session.New(), cfg)
	return svc
}
