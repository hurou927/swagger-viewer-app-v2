package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/swagger-viewer/swagger-viewer-app-v2/lib/db"
)

type ErrorElm struct {
	Code int `json:"code" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type ErrorBody struct {
	Error ErrorElm `json:"error" validate:"required"`
}


func CreateErrorResponse(statusCode int, errorBody ErrorBody)(events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(errorBody)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse {
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
	}
	return resp, nil
}

var a int = 4

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println(a)	
	serviceDao, err := servicedb.NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"));

	if err != nil {
		return CreateErrorResponse(500, ErrorBody {
			Error: ErrorElm {
				Code: 400,
				Message: "DynamoClientError",
			},
		})
	}


  	serviceDto, err := serviceDao.GetService(request.PathParameters["id"]);

	if err != nil || serviceDto == nil {
    	fmt.Println(err)
    	return events.APIGatewayProxyResponse {StatusCode: 404}, err
  	}
	

	body, err := json.Marshal(serviceDto)
	if err != nil {
		return events.APIGatewayProxyResponse {StatusCode: 404}, err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse {
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
	}

	return resp, nil
}

func main() {

	lambda.Start(Handler)
}