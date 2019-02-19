package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/swagger-viewer/swagger-viewer-app-v2/lib/common"
	versiondb "github.com/swagger-viewer/swagger-viewer-app-v2/lib/db/version"
)

var versionDao versiondb.VersionRepositoryDao
var versionInitError error

type requestBody struct {
	// ID      string `json:"id" validate:"required"`
	// Version string `json:"version" validate:"required"`
	Path string `json:"path" validate:"required"`
	// Lastupdated int64  `json:"lastupdated"` validate:"required"
	Enable bool   `json:"enable" validate:"required"`
	Tag    string `json:"tag" validate:"required"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// serviceDao, err := servicedb.NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"))

	// todo: duplicate ServiceName Check
	if versionInitError != nil {
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

	requestEntity := versiondb.VersionEntity{
		ID:          request.PathParameters["id"],
		Version:     request.PathParameters["version"],
		Path:        reqbody.Path,
		Lastupdated: time.Now().Unix() * 1000,
		Enable:      reqbody.Enable,
		Tag:         reqbody.Tag,
	}

	if _, err := versionDao.UpdateVersion(requestEntity); err != nil { //Todo: Error
		fmt.Println(err.(*common.Error).Error())
		if err.(*common.Error).Code == 1001 {
			return common.CreateErrorResponse(404, common.ErrorBody{
				Error: common.ErrorElm{
					Code:    10001,
					Message: "ID and version do not exists",
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

	resp, err := common.CreateResponse(200, requestEntity)

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
	versionDao, versionInitError = versiondb.NewDaoDefaultConfig(os.Getenv("VERSIONTABLENAME"))
	lambda.Start(Handler)
}
