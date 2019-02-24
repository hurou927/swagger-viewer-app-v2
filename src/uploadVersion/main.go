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
	Enable   bool   `json:"enable" validate:"required"`
	Tag      string `json:"tag" validate:"required"`
	Format   string `json:"format" validate:"required"`
	Contents string `json:"contents" validate:"required"`
	// Version  string `json:"version"` //optional
}

type swagger struct {
	Swagger string `json:"swagger" validate:"required"`
	Info    struct {
		Version string `json:"version" validate:"required"`
	} `json:"info" validate:"required"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	var fileFormat common.Format
	if reqbody.Format == "yaml" || reqbody.Format == "yml" {
		fileFormat = common.Yml
	} else if reqbody.Format == "json" {
		fileFormat = common.Json
	} else {
		return common.CreateErrorResponse(400, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1401,
				Message: "FileFormat Error",
			},
		})
	}

	swagger, err := common.ValidateSwagger(fileFormat, reqbody.Contents)
	if err != nil {
		return common.CreateErrorResponse(400, common.ErrorBody{
			Error: common.ErrorElm{
				Code:    1402,
				Message: "Swagger Error",
			},
		})
	}

	fmt.Printf("TODO: upload swagger: %+v\n", swagger)
	fmt.Println("=====================================")
	var ext string
	if fileFormat == common.Yml {
		ext = "json"
	} else {
		ext = "yml"
	}
	bucketName := os.Getenv("BUKCETNAME")
	keyName := fmt.Sprintf("swagger/%s/%s_%d.%s", request.PathParameters["id"], swagger.Info.Version, time.Now().Unix(), ext)

	requestEntity := versiondb.VersionEntity{
		ID:          request.PathParameters["id"],
		Version:     swagger.Info.Version,
		Path:        keyName,
		Lastupdated: time.Now().Unix() * 1000,
		Enable:      reqbody.Enable,
		Tag:         reqbody.Tag,
	}

	if _, err := versionDao.UploadVersion(requestEntity, bucketName, keyName, reqbody.Contents); err != nil { //Todo: Error
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

	resp, err := common.CreateResponse(200, reqbody)
	fmt.Println(resp)
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
