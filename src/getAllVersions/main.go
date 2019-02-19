package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/swagger-viewer/swagger-viewer-app-v2/lib/common"
	versiondb "github.com/swagger-viewer/swagger-viewer-app-v2/lib/db/version"
)

var versionDao versiondb.VersionRepositoryDao
var versionInitError error

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// serviceDao, err := servicedb.NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"))

	if versionInitError != nil {
		return common.CreateErrorResponse(500, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1500,
				Message: "DynamoClientError",
			},
		})
	}

	versions, err := versionDao.GetAllVersions(request.PathParameters["id"])

	if err != nil {
		fmt.Println(err)
		return common.CreateErrorResponse(500, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1500,
				Message: "DB Error",
			},
		})
	}

	if versions == nil {
		return common.CreateErrorResponse(404, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1404,
				Message: "Service Not Found",
			},
		})
	}
	resp, err := common.CreateResponse(200, map[string]interface{}{
		"Items": versions,
	})

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
