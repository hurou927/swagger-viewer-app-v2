package common

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorElm struct {
	Code    int    `json:"code" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type ErrorBody struct {
	Error ErrorElm `json:"error" validate:"required"`
}

func CreateErrorResponse(statusCode int, errorBody ErrorBody) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(errorBody)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                 "application/json; charset=utf-8",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
	}
	return resp, nil
}

func CreateResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	bodybytes, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, bodybytes)

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                 "application/json; charset=utf-8",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
	}
	return resp, nil
}
