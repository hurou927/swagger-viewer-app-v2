package servicedb

import (
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
	"github.com/swagger-viewer/swagger-viewer-app-v2/lib/common"
)

func ValidateServiceName(serviceName string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]*$`).MatchString(serviceName)
}

// ServiceEntity provides Service DB Record Contents
type ServiceEntity struct {
	Id            string `json:"id"`
	Servicename   string `json:"servicename"`
	Latestversion string `json:"latestversion"`
	Lastupdated   int64  `json:"lastupdated"`
}

// UpdateServiceEntity is used for UpdateServiceRepositoryDao
// if value is nil, it is not updated.
type UpdateServiceEntity struct {
	Id            *string `json:"id"`
	Servicename   *string `json:"servicename"`
	Latestversion *string `json:"latestversion"`
	Lastupdated   *int64  `json:"lastupdated"`
}

// ServiceRepositoryDao provides an interface of Dao for service db
// Mock:
// type serviceRepositoryDaoMock struct {
// 	ServiceRepositoryDao
// 	tableName string
// 	dynamoClient *dynamodb.DynamoDB
// }
//
// func (this *serviceRepositoryDaoMock) GetService(serviceId string) (*ServiceEntity, error) {
// 	return ServiceEntity{
// 		Id:            "mockid",
// 		Servicename:   "mockservice",
// 		Latestversion: "mock1.2.3",
// 		Lastupdated:   11111111,
// 	}, nil
// }
// .....
//

type ServiceRepositoryDao interface {
	GetService(serviceId string) (*ServiceEntity, error)
	GetServiceList() ([]ServiceEntity, error)
	CreateService(service ServiceEntity) (*ServiceEntity, error)
	UpdateService(service UpdateServiceEntity) (*ServiceEntity, error)
	DeleteService(serviceId string) (*ServiceEntity, error)
}

type serviceRepositoryDaoImpl struct {
	tableName    string
	dynamoClient *dynamodb.DynamoDB
}

// NewDaoDefaultConfig return DynamoDB Session
func NewDaoDefaultConfig(tableName string) (ServiceRepositoryDao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.DisableEndpointHostPrefix = true

	if err != nil {
		return nil, common.NewError(200, "aws-sdk config error", err)
	}

	return &serviceRepositoryDaoImpl{
		dynamoClient: dynamodb.New(cfg),
		tableName:    tableName,
	}, nil
}

// NewDaoWithRegion return DynamoDB Session
func NewDaoWithRegion(tableName string, region string) (ServiceRepositoryDao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.Region = region
	cfg.DisableEndpointHostPrefix = true
	if err != nil {
		return nil, common.NewError(200, "aws-sdk config error", err)
	}
	return &serviceRepositoryDaoImpl{
		dynamoClient: dynamodb.New(cfg),
		tableName:    tableName,
	}, nil
}

// NewDaoWithRegionAndEndpoint return DynamoDB Session
// If you are using dynamodb local, use it.
// example: dao, err := NewDaoWithRegionAndEndpoint("tablename", "ap-northeast-1", "http://localhost:8000")
func NewDaoWithRegionAndEndpoint(tableName string, region string, endpoint string) (ServiceRepositoryDao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(endpoint)
	cfg.Region = region
	cfg.DisableEndpointHostPrefix = true
	if err != nil {
		return nil, common.NewError(200, "aws-sdk config error", err)
	}

	return &serviceRepositoryDaoImpl{
		dynamoClient: dynamodb.New(cfg),
		tableName:    tableName,
	}, nil
}

// GetService gets a service info.
func (this *serviceRepositoryDaoImpl) GetService(serviceId string) (*ServiceEntity, error) {
	if this == nil {
		return nil, common.NewError(100, "nil pointer receiver", nil)
	}
	result, err := this.dynamoClient.GetItemRequest(&dynamodb.GetItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"id": {
				S: aws.String(serviceId),
			},
		},
		TableName: aws.String(this.tableName),
	}).Send()

	if err != nil {
		return nil, common.NewError(300, "dynamoDB error", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	entity := ServiceEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &entity); err != nil {
		return nil, common.NewError(101, "unmarshal error", err)
	}
	return &entity, nil
}

// GetServiceList get all service stored in db
func (this *serviceRepositoryDaoImpl) GetServiceList() ([]ServiceEntity, error) {
	if this == nil {
		return nil, common.NewError(100, "nil pointer receiver", nil)
	}
	req := this.dynamoClient.ScanRequest(&dynamodb.ScanInput{
		TableName: aws.String(this.tableName),
	})
	p := req.Paginate()

	var items []map[string]dynamodb.AttributeValue
	for p.Next() {
		page := p.CurrentPage()
		items = append(items, page.Items...)
	}
	if err := p.Err(); err != nil {
		return nil, common.NewError(300, "dynamodb scan paginate error", err)
	}

	var services []ServiceEntity
	if err := dynamodbattribute.UnmarshalListOfMaps(items, &services); err != nil {
		return nil, common.NewError(301, "dynamoDB unmarhsallist error", err)
	}
	return services, nil
}

// CreateService creates service
func (this *serviceRepositoryDaoImpl) CreateService(service ServiceEntity) (*ServiceEntity, error) {
	if this == nil {
		return nil, common.NewError(100, "nil pointer receiver", nil)
	}

	item, err := dynamodbattribute.MarshalMap(service)
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
	entity := ServiceEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &entity); err != nil {
		return nil, common.NewError(301, "dynamoDB unmarhsallist error", err)
	}
	return &entity, nil // return old data. Usually, This value is nothing.
}

// UpdateService updates service info.
func (this *serviceRepositoryDaoImpl) UpdateService(service UpdateServiceEntity) (*ServiceEntity, error) {
	if this == nil {
		return nil, common.NewError(100, "nil pointer receiver", nil)
	}
	if service.Id == nil {
		return nil, common.NewError(1001, "id is required", nil)
	}

	var willBeUpdated bool = false

	var update expression.UpdateBuilder = expression.UpdateBuilder{}
	if service.Servicename != nil {
		willBeUpdated = true
		update = update.Set(expression.Name("servicename"), expression.Value(*service.Servicename))
	}
	if service.Latestversion != nil {
		willBeUpdated = true
		update = update.Set(expression.Name("latestversion"), expression.Value(*service.Latestversion))
	}
	if service.Lastupdated != nil {
		willBeUpdated = true
		update = update.Set(expression.Name("lastupdated"), expression.Value(*service.Lastupdated))
	}

	if !willBeUpdated {
		return nil, common.NewError(1001, "one or more attributes are required", nil)
	}

	condition := expression.AttributeExists(expression.Name("id"))
	// anotherCondition := expression.Not(condition)

	expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(condition).Build()
	if err != nil {
		return nil, common.NewError(302, "expression build error", err)
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String(this.tableName),
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		Key: map[string]dynamodb.AttributeValue{
			"id": {
				S: aws.String(*service.Id),
			},
		},
		ReturnValues: dynamodb.ReturnValueAllNew,
	}
	// fmt.Printf("%+v\n", input)
	result, err := this.dynamoClient.UpdateItemRequest(input).Send()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return nil, common.NewError(1002, "id does not exists", aerr)
			default:
				return nil, common.NewError(300, "dynamodb put error", aerr)
			}
		} else {
			return nil, common.NewError(0, "unknown error", err)
		}
	}

	entity := ServiceEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &entity); err != nil {
		return nil, common.NewError(301, "dynamoDB unmarhsallist error", err)
	}
	return &entity, nil
}

// DeleteService deletes service info
func (this *serviceRepositoryDaoImpl) DeleteService(serviceId string) (*ServiceEntity, error) {
	if this == nil {
		return nil, common.NewError(100, "nil pointer receiver", nil)
	}
	result, err := this.dynamoClient.DeleteItemRequest(&dynamodb.DeleteItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"id": {
				S: aws.String(serviceId),
			},
		},
		TableName:    aws.String(this.tableName),
		ReturnValues: dynamodb.ReturnValueAllOld,
	}).Send()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, common.NewError(300, "dynamodb put error", aerr)
		}

		return nil, common.NewError(0, "unknown error", err)
	}

	entity := ServiceEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &entity); err != nil {
		return nil, common.NewError(301, "dynamoDB unmarhsallist error", err)
	}
	return &entity, nil
}
