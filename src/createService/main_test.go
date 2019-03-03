package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/swagger-viewer/swagger-viewer-app-v2/lib/common"
	servicedb "github.com/swagger-viewer/swagger-viewer-app-v2/lib/db"
)

var dynamoLocalEndpoint string = "http://localhost:8027"

func TestHandlerSuccess(t *testing.T) {

	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	serviceDao, serviceInitError = servicedb.NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), dynamoLocalEndpoint)

	body := map[string]interface{}{
		"serviceName": "service",
	}
	queryParams := map[string]string{}
	pathParams := map[string]string{}

	request, err := common.CreateProxyRequest(body, queryParams, pathParams)

	var ctx context.Context
	response, err := Handler(ctx, request)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if response.StatusCode != 201 {
		fmt.Printf("%+v\n", response)
		t.Fatalf("error response %d", response.StatusCode)
	}
	fmt.Printf("%+v\n", response.Body)
}
