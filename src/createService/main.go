package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/swagger-viewer/swagger-viewer-app-v2/lib/common"
	servicedb "github.com/swagger-viewer/swagger-viewer-app-v2/lib/db"
)

var serviceDao servicedb.ServiceRepositoryDao
var serviceInitError error

type requestBody struct {
	Servicename   string `json:"servicename" validate:"required"`
	Latestversion string `json:"latestversion" validate:"required"`
	Lastupdated   int64  `json:"lastupdated" validate:"required"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// serviceDao, err := servicedb.NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"))

	// todo: dupulicate ServiceName Check
	if serviceInitError != nil {
		return common.CreateErrorResponse(500, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1500,
				Message: "DynamoClientError",
			},
		})
	}

	var reqbody requestBody
	if err := json.Unmarshal([]byte(request.Body), &reqbody); err != nil {
		return common.CreateErrorResponse(500, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1500,
				Message: "Internal Error",
			},
		})
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return common.CreateErrorResponse(500, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1500,
				Message: "Internal Error",
			},
		})
	}

	requestEntity := servicedb.ServiceEntity{
		Id:            id.String(),
		Servicename:   reqbody.Servicename,
		Latestversion: reqbody.Latestversion,
		Lastupdated:   reqbody.Lastupdated,
	}

	if _, err := serviceDao.CreateService(requestEntity); err != nil { //Todo: Error
		return common.CreateErrorResponse(400, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1400,
				Message: "DynamoError",
			},
		})
	}

	body, err := json.Marshal(map[string]interface{}{
		"item": requestEntity,
	})
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
