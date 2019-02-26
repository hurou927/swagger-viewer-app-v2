package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/swagger-viewer/swagger-viewer-app-v2/lib/common"
	versiondb "github.com/swagger-viewer/swagger-viewer-app-v2/lib/db/version"
)

var dynamoLocalEndpoint string = "http://localhost:8027"
var s3LocalEndpoint string = "http://localhost:4568"

func TestHandlerSuccess(t *testing.T) {

	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("VERSIONTABLENAME", "swagger-dev-swagger-dynamo-versioninfo")
	os.Setenv("SWAGGER_BUCKET_NAME", "swagger-repository-test")

	versionDao, versionInitError = versiondb.NewDaoWithEndpoints(
		os.Getenv("VERSIONTABLENAME"),
		versiondb.AwsEndpoint{
			Region:   os.Getenv("AWS_DEFAULT_REGION"),
			Endpoint: dynamoLocalEndpoint,
		},
		versiondb.AwsEndpoint{
			Region:   os.Getenv("AWS_DEFAULT_REGION"),
			Endpoint: s3LocalEndpoint,
		},
	)

	yamlInput := `
swagger: '2.0'
info:
  description: これはアパートに関するAPIです。
  version: 0.0.1
  title: アパートAPI
`

	body := map[string]interface{}{
		"enable":   true,
		"Contents": yamlInput,
		"Format":   "yaml",
		"Version":  "1.2.1",
		"tag":      "nonono",
	}
	queryParams := map[string]string{}
	pathParams := map[string]string{
		"id": "524f25fe-b711-3ae8-b7b8-93fffaaeb4e0",
	}

	request, err := common.CreateProxyRequest(body, queryParams, pathParams)

	var ctx context.Context
	response, err := Handler(ctx, request)
	fmt.Println(response, err)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("error response %d", response.StatusCode)
	}

}
