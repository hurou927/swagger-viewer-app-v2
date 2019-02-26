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

func TestHandlerSuccess(t *testing.T) {

	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("VERSIONTABLENAME", "swagger-dev-swagger-dynamo-versioninfo")

	versionDao, versionInitError = versiondb.NewDaoWithRegionAndEndpoint(os.Getenv("VERSIONTABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), dynamoLocalEndpoint)

	body := map[string]interface{}{}
	queryParams := map[string]string{}
	pathParams := map[string]string{
		"id": "524f25fe-b711-3ae8-b7b8-93fffaaeb4e0",
	}

	request, err := common.CreateProxyRequest(body, queryParams, pathParams)

	var ctx context.Context
	response, err := Handler(ctx, request)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("error response %d", response.StatusCode)
	}
	fmt.Printf("%+v\n", response.Body)
}
