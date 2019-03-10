package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
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

	// serviceDao, err := servicedb.NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"))

	// todo: duplicate ServiceName Check
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
		Latestversion: "0.0.0",
		Lastupdated:   time.Now().Unix() * 1000,
	}

	if _, err := serviceDao.CreateService(requestEntity); err != nil { //Todo: Error
		if err.(*common.Error).Code == 1001 {
			return common.CreateErrorResponse(400, common.ErrorBody{
				Error: common.ErrorElm{
					Code:    10001,
					Message: "ID already exist",
				},
			})
		}
		return common.CreateErrorResponse(400, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1400,
				Message: "DynamoError",
			},
		})
	}

	resp, err := common.CreateResponse(201, requestEntity)

	if err != nil {
		return common.CreateErrorResponse(500, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1500,
				Message: "Internal Error",
			},
		})
	}

	return resp, nil
}

func main() {
	serviceDao, serviceInitError = servicedb.NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"))
	lambda.Start(Handler)
}
