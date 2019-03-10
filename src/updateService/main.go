package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/swagger-viewer/swagger-viewer-app-v2/lib/common"
	servicedb "github.com/swagger-viewer/swagger-viewer-app-v2/lib/db"
)

var serviceDao servicedb.ServiceRepositoryDao
var serviceInitError error

type requestBody struct {
	Servicename string `json:"servicename" validate:"required"`
	// Latestversion string `json:"latestversion" validate:"required"`
	// Lastupdated   int64  `json:"lastupdated" validate:"required"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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
	if servicedb.ValidateServiceName(reqbody.Servicename) == false {
		return common.CreateErrorResponse(400, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1301,
				Message: "Service Name must be url-safe(^[a-zA-Z0-9_-]*$)",
			},
		})
	}

	updateService := servicedb.UpdateServiceEntity{}
	var serviceId = request.PathParameters["id"]
	updateService.Id = &serviceId
	updateService.Servicename = &reqbody.Servicename

	// if reqbody.Servicename != "" {
	// 	updateService.Servicename = &reqbody.Servicename
	// }
	// if reqbody.Lastupdated != 0 {
	// 	updateService.Lastupdated = &reqbody.Lastupdated
	// }
	// if reqbody.Latestversion != "" {
	// 	updateService.Latestversion = &reqbody.Latestversion
	// }

	if _, err := serviceDao.UpdateService(updateService); err != nil {
		if err.(*common.Error).Code == 1002 {
			return common.CreateErrorResponse(404, common.ErrorBody{
				Error: common.ErrorElm{
					Code:    10002,
					Message: "ID does not exist",
				},
			})
		}
		return common.CreateErrorResponse(400, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1401,
				Message: "Internal Error",
			},
		})
	}

	body, err := json.Marshal(map[string]interface{}{
		"success": true,
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
