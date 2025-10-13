package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rendyfutsuy/base-go.git/utils"
)

func NewSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(utils.ConfigVars.String("aws.aws_endpoint")),
		Region:   aws.String(utils.ConfigVars.String("aws.aws_region")),
		Credentials: credentials.NewStaticCredentials(
			utils.ConfigVars.String("aws.aws_access_key"),
			utils.ConfigVars.String("aws.aws_access_secret"),
			"",
		),
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}

func StartS3() *s3.S3 {
	sess, err := NewSession()
	if err != nil {
		fmt.Println("Failed to create AWS session:", err)
		return nil
	}

	s3Client := s3.New(sess)
	utils.Logger.Info("S3 session & client initialized")

	return s3Client
}
