package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/swagger-viewer/swagger-viewer-app-v2/lib/common"
	servicedb "github.com/swagger-viewer/swagger-viewer-app-v2/lib/db"
)

var serviceDao servicedb.ServiceRepositoryDao
var serviceInitError error

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// serviceDao, err := servicedb.NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"))

	if serviceInitError != nil {
		return common.CreateErrorResponse(500, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1500,
				Message: "DynamoClientError",
			},
		})
	}

	serviceEntity, err := serviceDao.GetService(request.PathParameters["id"])

	if err != nil {
		fmt.Println(err)
		return common.CreateErrorResponse(500, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1500,
				Message: "DB Error",
			},
		})
	}

	if serviceEntity == nil {
		return common.CreateErrorResponse(404, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1404,
				Message: "Service Not Found",
			},
		})
	}

	body, err := json.Marshal(serviceEntity)
	if err != nil {
		return common.CreateErrorResponse(500, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1500,
				Message: "Internal Error",
			},
		})
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
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
	serviceDao, serviceInitError = servicedb.NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"))
	lambda.Start(Handler)
}
