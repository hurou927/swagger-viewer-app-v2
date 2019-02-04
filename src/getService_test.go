package main

import (
    "fmt"
    "testing"
    "os"
    "github.com/swagger-viewer/swagger-viewer-app-v2/src/lib"
    "context"
)



func TestHandlerSuccess(t *testing.T) {
    os.Setenv("AWS_DEFAULT_REGION","ap-northeast-1")
    os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")
    body := map[string] interface{} {}
	queryParams := map[string] string {}
	pathParams := map[string] string {
		"id": "524f25fe-b711-3ae8-b7b8-93fffaaeb4e0",
    }
    

    request, err := request.CreateProxyRequest(body,queryParams, pathParams)
    
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


func TestHandlerFailure(t *testing.T) {
    os.Setenv("AWS_DEFAULT_REGION","ap-northeast-1")
    os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")
    body := map[string] interface{} {}
	queryParams := map[string] string {}
	pathParams := map[string] string {
		"id": "524f25",
    }
    

    request, err := request.CreateProxyRequest(body,queryParams, pathParams)
    
    var ctx context.Context
    response, err := Handler(ctx, request)

    if err != nil {
        t.Fatalf("failed test %#v", err)
    }
    if response.StatusCode != 404 {
        t.Fatalf("error response %d", response.StatusCode)
    }
    // fmt.Printf("%+v\n", response.Body)
}