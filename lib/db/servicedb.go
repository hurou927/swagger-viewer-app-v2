package servicedb

import (
  "fmt"
  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/aws/external"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)


type Dto struct {
  Id			string `json:"id"`
  Servicename	string `json:"servicename"`
  Latestversion	string `json:"latestversion"`
  Lastupdated	int64  `json:"lastupdated"`
}

type Dao struct {
	tableName string
	dynamoClient *dynamodb.DynamoDB
}





func NewDaoDefaultConfig(tableName string)(*Dao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
  cfg.DisableEndpointHostPrefix = true
	if err != nil {
		return &Dao{}, err
	}

  	return &Dao{
			dynamoClient: dynamodb.New(cfg),
			tableName: tableName,
	}, nil
}

func NewDaoWithRegion(tableName string, region string)(*Dao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
  	cfg.Region = region
  	cfg.DisableEndpointHostPrefix = true
  	if err != nil {
		return &Dao{}, err
	}

  	return &Dao{
			dynamoClient: dynamodb.New(cfg),
			tableName: tableName,
	}, nil
}

func NewDaoWithRegionAndEndpoint(tableName string, region string, endpoint string)(*Dao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(endpoint)
  	cfg.Region = region
  	cfg.DisableEndpointHostPrefix = true
  	if err != nil {
		return &Dao{}, err
	}

  	return &Dao{
			dynamoClient: dynamodb.New(cfg),
			tableName: tableName,
	}, nil
}



func (this *Dao) GetService (servicId string) (*Dto, error) {
  if this == nil {
    return nil, fmt.Errorf("nil pointer receiver")
  }
	result, err  := this.dynamoClient.GetItemRequest(&dynamodb.GetItemInput{
	  Key: map[string]dynamodb.AttributeValue{
	    "id": {
	      S: aws.String(servicId),
	    },
	  },
	  TableName: aws.String(this.tableName),
	}).Send()
  
  if err != nil{
    return nil, err
  }

  if result.Item == nil {
    return nil, nil
  }


  dto := Dto{}
  if err := dynamodbattribute.UnmarshalMap(result.Item, &dto); err!= nil{
    return &Dto{}, err
  }
  return &dto, nil
}



func (this *Dao) GetServiceList () ([]Dto, error) {
  if this == nil {
    return nil, fmt.Errorf("nil pointer receiver")
  }
  req := this.dynamoClient.ScanRequest(&dynamodb.ScanInput{
    TableName: aws.String(this.tableName),
  })
  p := req.Paginate()

	var items []map[string] dynamodb.AttributeValue
  for p.Next() {
		page := p.CurrentPage()
    items = append(items, page.Items ...)
	}
  if err := p.Err(); err != nil {
		return nil, err
  }

	var servicies []Dto
	if err := dynamodbattribute.UnmarshalListOfMaps(items, &servicies); err!= nil{
    	return nil, err
  }

	return servicies, nil

}



func (this *Dao) PostServices (services []Dto) ([]Dto, error) {
  if this == nil {
    return nil, fmt.Errorf("nil pointer receiver")
  }

  UnprocessedItems := make([]Dto, len(services), len(services))

  for index := range services {
    item, err := dynamodbattribute.MarshalMap(services[index])
    if err != nil {
      fmt.Println(err)
      UnprocessedItems[0] = services[index]
      continue;
    } 
    result, err := this.dynamoClient.PutItemRequest(&dynamodb.PutItemInput{
      TableName: aws.String(this.tableName),
      Item: item,
      ConditionExpression: aws.String("attribute_not_exists(#userid)"),
      ExpressionAttributeNames: map[string] string  {
        "#id" : "id",
      },
      ReturnConsumedCapacity: dynamodb.ReturnConsumedCapacityTotal,
    }).Send()

    if err != nil {
      fmt.Println(err)
      UnprocessedItems[0] = services[index]
      continue;
    }

    fmt.Println(result)
  }
  return UnprocessedItems, nil
}


// func (this *Dao) UpdateService (servicie Dto) ()