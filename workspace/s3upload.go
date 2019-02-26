package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3LocalEndpoint string = "http://localhost:4568"

func main() {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.DisableEndpointHostPrefix = true
	if err != nil {
		fmt.Println(err)
		return
	}

	// cfg, err := external.LoadDefaultAWSConfig()
	// cfg.EndpointResolver = aws.ResolveWithEndpointURL(s3LocalEndpoint)
	// cfg.Region = os.Getenv("AWS_DEFAULT_REGION")
	// cfg.DisableEndpointHostPrefix = true
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// s3Uploader := s3manager.NewUploader(cfg)

	yamlInput := `
swagger: '2.0'
info:
  description: これはアパートに関するAPIです。
  version: 0.0.1
  title: アパートAPI
`

	bucket := "swagger-repository-test"
	key := "src/test.yml"

	svc := s3.New(cfg)
	svc.ForcePathStyle = true
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(yamlInput)),
	}

	result, err := svc.PutObjectRequest(input).Send()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Success: %+v\n", result)
}
