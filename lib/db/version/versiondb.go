package versiondb

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/swagger-viewer/swagger-viewer-app-v2/lib/common"
)

type VersionEntity struct {
	ID          string `json:"id"`
	Version     string `json:"version"`
	Path        string `json:"path"`
	Lastupdated int64  `json:"lastupdated"`
	Enable      bool   `json:"enable"`
	Tag         string `json:"tag"`
}

type UpdateVersionEntity struct {
	ID          *string `json:"id"`
	Version     *string `json:"version"`
	Path        *string `json:"path"`
	Lastupdated *int64  `json:"lastupdated"`
	Enable      *bool   `json:"enable"`
	Tag         *string `json:"tag"`
}

type AwsEndpoint struct {
	Region   string
	Endpoint string
}

type VersionRepositoryDao interface {
	GetAllVersions(servicId string) ([]VersionEntity, error)
	CreateVersion(version VersionEntity) (*VersionEntity, error)
	UpdateVersion(version VersionEntity) (*VersionEntity, error)
	UploadVersion(version VersionEntity, bucket string, key string, contents string) (*VersionEntity, error)
}

type versionRepositoryDaoImpl struct {
	tableName    string
	dynamoClient *dynamodb.DynamoDB
	s3Client     *s3.S3
}

// NewDaoDefaultConfig return DynamoDB Session
func NewDaoDefaultConfig(tableName string) (VersionRepositoryDao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.DisableEndpointHostPrefix = true

	if err != nil {
		return nil, common.NewError(200, "aws-sdk config error", err)
	}

	s3Client := s3.New(cfg)
	s3Client.ForcePathStyle = true

	return &versionRepositoryDaoImpl{
		dynamoClient: dynamodb.New(cfg),
		tableName:    tableName,
		s3Client:     s3Client,
	}, nil
}

// NewDaoWithRegion return DynamoDB Session
func NewDaoWithRegion(tableName string, region string) (VersionRepositoryDao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.Region = region
	cfg.DisableEndpointHostPrefix = true
	if err != nil {
		return nil, common.NewError(200, "aws-sdk config error", err)
	}
	s3Client := s3.New(cfg)
	// s3Client.ForcePathStyle = true

	return &versionRepositoryDaoImpl{
		dynamoClient: dynamodb.New(cfg),
		tableName:    tableName,
		s3Client:     s3Client,
	}, nil
}

// NewDaoWithRegionAndEndpoint return DynamoDB Session
// If you are using dynamodb local, use it.
// example: dao, err := NewDaoWithRegionAndEndpoint("tablename", "ap-northeast-1", "http://localhost:8000")
func NewDaoWithRegionAndEndpoint(tableName string, dynamoRegion string, dynamoEndpoint string) (VersionRepositoryDao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(dynamoEndpoint)
	cfg.Region = dynamoRegion
	cfg.DisableEndpointHostPrefix = true
	if err != nil {
		return nil, common.NewError(200, "aws-sdk config error", err)
	}

	s3Client := s3.New(cfg)
	s3Client.ForcePathStyle = true

	return &versionRepositoryDaoImpl{
		dynamoClient: dynamodb.New(cfg),
		tableName:    tableName,
		s3Client:     s3Client,
	}, nil
}

func NewDaoWithEndpoints(tableName string, dynamoEndpoint AwsEndpoint, s3Endpoint AwsEndpoint) (VersionRepositoryDao, error) {
	dynamoCfg, err := external.LoadDefaultAWSConfig()
	dynamoCfg.EndpointResolver = aws.ResolveWithEndpointURL(dynamoEndpoint.Endpoint)
	dynamoCfg.Region = dynamoEndpoint.Region
	dynamoCfg.DisableEndpointHostPrefix = true
	if err != nil {
		return nil, common.NewError(200, "aws-sdk config error", err)
	}

	s3Cfg, err := external.LoadDefaultAWSConfig()
	s3Cfg.EndpointResolver = aws.ResolveWithEndpointURL(s3Endpoint.Endpoint)
	s3Cfg.Region = s3Endpoint.Region
	s3Cfg.DisableEndpointHostPrefix = true
	if err != nil {
		return nil, common.NewError(200, "aws-sdk config error", err)
	}

	s3Client := s3.New(s3Cfg)
	s3Client.ForcePathStyle = true

	return &versionRepositoryDaoImpl{
		dynamoClient: dynamodb.New(dynamoCfg),
		tableName:    tableName,
		s3Client:     s3Client,
	}, nil
}

func (this *versionRepositoryDaoImpl) GetAllVersions(serviceId string) ([]VersionEntity, error) {
	if this == nil {
		return nil, common.NewError(100, "nil pointer receiver", nil)
	}
	// keyCond := expression.Key("serviceid").Equal(expression.Value(serviceId))
	// expression, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	keyCond := expression.Key("id").Equal(expression.Value(serviceId))
	// proj := expression.NamesList(expression.Name("aName"), expression.Name("anotherName"), expression.Name("oneOtherName"))
	// builder := expression.NewBuilder().WithKeyCondition(keyCond).WithProjection(proj)
	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	expression, err := builder.Build()
	if err != nil {
		return nil, common.NewError(302, "expression build error", err)
	}

	queryInput := &dynamodb.QueryInput{
		KeyConditionExpression: expression.KeyCondition(),
		// ProjectionExpression:      expression.Projection(),
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		TableName:                 aws.String(this.tableName),
	}

	// fmt.Println(queryInput)
	result, err := this.dynamoClient.QueryRequest(queryInput).Send()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return nil, common.NewError(1000, "id already exists", aerr)
			default:
				return nil, common.NewError(300, "dynamodb query error", aerr)
			}
		}
		return nil, common.NewError(0, "unknown error", err)
	}

	if result.Items == nil {
		return nil, nil
	}

	var versions []VersionEntity
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &versions); err != nil {
		return nil, common.NewError(301, "dynamoDB unmarhsallist error", err)
	}
	return versions, nil
}

func (this *versionRepositoryDaoImpl) CreateVersion(version VersionEntity) (*VersionEntity, error) {
	if this == nil {
		return nil, common.NewError(100, "nil pointer receiver", nil)
	}

	item, err := dynamodbattribute.MarshalMap(version)
	if err != nil {
		return nil, common.NewError(301, "dynamoDB marhsallist error", err)
	}
	result, err := this.dynamoClient.PutItemRequest(&dynamodb.PutItemInput{
		TableName:           aws.String(this.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(#id)"),
		ExpressionAttributeNames: map[string]string{
			"#id": "id",
		},
		ReturnConsumedCapacity: dynamodb.ReturnConsumedCapacityTotal,
	}).Send()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return nil, common.NewError(1000, "id already exists", aerr)
			default:
				return nil, common.NewError(300, "dynamodb put error", aerr)
			}
		}
		return nil, common.NewError(0, "unknown error", err)
	}
	entity := VersionEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &entity); err != nil {
		return nil, common.NewError(301, "dynamoDB unmarhsallist error", err)
	}
	return &entity, nil // return old data. Usually, This value is nothing.
}

func (this *versionRepositoryDaoImpl) UpdateVersion(version VersionEntity) (*VersionEntity, error) {
	if this == nil {
		return nil, common.NewError(100, "nil pointer receiver", nil)
	}

	item, err := dynamodbattribute.MarshalMap(version)
	if err != nil {
		return nil, common.NewError(301, "dynamoDB marhsallist error", err)
	}
	result, err := this.dynamoClient.PutItemRequest(&dynamodb.PutItemInput{
		TableName:           aws.String(this.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_exists(#id) AND attribute_exists(#version)"),
		ExpressionAttributeNames: map[string]string{
			"#id":      "id",
			"#version": "version",
		},
		ReturnConsumedCapacity: dynamodb.ReturnConsumedCapacityTotal,
	}).Send()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return nil, common.NewError(1001, "id and version do not exist", aerr)
			default:
				return nil, common.NewError(300, "dynamodb put error", aerr)
			}
		}
		return nil, common.NewError(0, "unknown error", err)
	}
	entity := VersionEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &entity); err != nil {
		return nil, common.NewError(301, "dynamoDB unmarhsallist error", err)
	}
	return &entity, nil // return old data. Usually, This value is nothing.
}

func (this *versionRepositoryDaoImpl) UploadVersion(version VersionEntity, bucket string, key string, contents string) (*VersionEntity, error) {
	if this == nil {
		return nil, common.NewError(100, "nil pointer receiver", nil)
	}

	item, err := dynamodbattribute.MarshalMap(version)
	if err != nil {
		return nil, common.NewError(301, "dynamoDB marhsallist error", err)
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(contents)),
	}

	s3Result, err := this.s3Client.PutObjectRequest(input).Send()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil, common.NewError(301, "s3 putobject error", err)
	}

	fmt.Printf("%+v\n", s3Result)

	result, err := this.dynamoClient.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(this.tableName),
		Item:      item,
		ExpressionAttributeNames: map[string]string{
			"#id":      "id",
			"#version": "version",
		},
		ReturnConsumedCapacity: dynamodb.ReturnConsumedCapacityTotal,
	}).Send()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return nil, common.NewError(1001, "id and version do not exist", aerr)
			default:
				return nil, common.NewError(300, "dynamodb put error", aerr)
			}
		}
		return nil, common.NewError(0, "unknown error", err)
	}
	entity := VersionEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &entity); err != nil {
		return nil, common.NewError(301, "dynamoDB unmarhsallist error", err)
	}
	return &entity, nil // return old data. Usually, This value is nothing.
}
