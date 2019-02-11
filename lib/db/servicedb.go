package servicedb

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
)

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
// func (this *serviceRepositoryDaoMock) GetService(servicId string) (*ServiceEntity, error) {
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
	GetService(servicId string) (*ServiceEntity, error)
	GetServiceList() ([]ServiceEntity, error)
	CreateService(service ServiceEntity) (*ServiceEntity, error)
	UpdateService(service UpdateServiceEntity) (*ServiceEntity, error)
	DeleteService(servicId string) (*ServiceEntity, error)
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
		return nil, err
	}

	return &serviceRepositoryDaoImpl{
		dynamoClient: dynamodb.New(cfg),
		tableName:    tableName,
	}, nil
}

// NewDaoDefaultConfig return DynamoDB Session
func NewDaoWithRegion(tableName string, region string) (ServiceRepositoryDao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.Region = region
	cfg.DisableEndpointHostPrefix = true
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &serviceRepositoryDaoImpl{
		dynamoClient: dynamodb.New(cfg),
		tableName:    tableName,
	}, nil
}

// NewDaoDefaultConfig return DynamoDB Session
// If you are using dynamodb local, use it.
// example: dao, err := NewDaoWithRegionAndEndpoint("tablename", "ap-northeast-1", "http://localhost:8000")
func NewDaoWithRegionAndEndpoint(tableName string, region string, endpoint string) (ServiceRepositoryDao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(endpoint)
	cfg.Region = region
	cfg.DisableEndpointHostPrefix = true
	if err != nil {
		return nil, err
	}

	return &serviceRepositoryDaoImpl{
		dynamoClient: dynamodb.New(cfg),
		tableName:    tableName,
	}, nil
}

func (this *serviceRepositoryDaoImpl) GetService(servicId string) (*ServiceEntity, error) {
	if this == nil {
		return nil, fmt.Errorf("nil pointer receiver")
	}
	result, err := this.dynamoClient.GetItemRequest(&dynamodb.GetItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"id": {
				S: aws.String(servicId),
			},
		},
		TableName: aws.String(this.tableName),
	}).Send()

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	entity := ServiceEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &entity); err != nil {
		return &ServiceEntity{}, err
	}
	return &entity, nil
}

func (this *serviceRepositoryDaoImpl) GetServiceList() ([]ServiceEntity, error) {
	if this == nil {
		return nil, fmt.Errorf("nil pointer receiver")
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
		return nil, err
	}

	var servicies []ServiceEntity
	if err := dynamodbattribute.UnmarshalListOfMaps(items, &servicies); err != nil {
		return nil, err
	}

	return servicies, nil

}

func (this *serviceRepositoryDaoImpl) CreateService(service ServiceEntity) (*ServiceEntity, error) {
	if this == nil {
		return nil, fmt.Errorf("nil pointer receiver")
	}

	item, err := dynamodbattribute.MarshalMap(service)
	if err != nil {
		return nil, err
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
			return nil, aerr
		}
		return nil, err
	}
	entity := ServiceEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &entity); err != nil {
		return &ServiceEntity{}, err
	}
	return &entity, nil // return old data. Usually, This value is nothing.
}

func (this *serviceRepositoryDaoImpl) UpdateService(service UpdateServiceEntity) (*ServiceEntity, error) {
	if this == nil {
		return nil, fmt.Errorf("nil pointer receiver")
	}
	if service.Id == nil {
		return nil, fmt.Errorf("ID is required")
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
		return nil, fmt.Errorf("not updated")
	}

	condition := expression.AttributeExists(expression.Name("id"))
	// anotherCondition := expression.Not(condition)

	expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(condition).Build()
	if err != nil {
		return nil, err
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
				return nil, fmt.Errorf("dynamodb:ErrCodeConditionalCheckFailedException:Id already exists")
			default:
				return nil, aerr
			}
		} else {
			return nil, err
		}
	}

	entity := ServiceEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (this *serviceRepositoryDaoImpl) DeleteService(servicId string) (*ServiceEntity, error) {
	if this == nil {
		return nil, fmt.Errorf("nil pointer receiver")
	}
	result, err := this.dynamoClient.DeleteItemRequest(&dynamodb.DeleteItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"id": {
				S: aws.String(servicId),
			},
		},
		TableName:    aws.String(this.tableName),
		ReturnValues: dynamodb.ReturnValueAllOld,
	}).Send()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, aerr
		}

		return nil, err
	}

	entity := ServiceEntity{}
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &entity); err != nil {
		return &ServiceEntity{}, err
	}
	return &entity, nil
}
