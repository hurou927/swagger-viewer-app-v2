package servicedb

import (
	// "fmt"

	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

// var dynamoLocalEndpoint string = "http://localhost:8027"

// post -> get -> delete -> get

func TestPOST_GET(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), dynamoLocalEndpoint)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	id, _ := uuid.NewUUID()
	serviceId := id.String()
	getService0, err := dao.GetService(serviceId)
	if err != nil || getService0 != nil {
		t.Fatalf("failed test %#v", err)
	}
	createdService := ServiceEntity{
		Id:            serviceId,
		Servicename:   "testservice",
		Latestversion: "1.2.3",
		Lastupdated:   53,
	}
	oldService, err := dao.CreateService(createdService)
	if err != nil || oldService.Id != "" {
		t.Fatalf("failed test(need to initialize dynamodb local) %#v", err)
	}
	getService1, err := dao.GetService(serviceId)
	if err != nil || getService1 == nil {
		t.Fatalf("no item(create service error) %#v", err)
	}
	if diff := cmp.Diff(*getService1, createdService); diff != "" {
		t.Fatalf("failed test(created service is wrong) %#v", err)
	}
}

func TestPOST_UPDATE_GET(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), dynamoLocalEndpoint)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	id, _ := uuid.NewUUID()
	serviceId := id.String()
	getService0, err := dao.GetService(serviceId)
	if err != nil || getService0 != nil {
		t.Fatalf("failed test %#v", err)
	}
	createdService := ServiceEntity{
		Id:            serviceId,
		Servicename:   "testservice",
		Latestversion: "1.2.3",
		Lastupdated:   53,
	}

	if _, err := dao.CreateService(createdService); err != nil {
		t.Fatalf("failed test(need to initialize dynamodb local) %#v", err)
	}

	servicename := "updateServiceErrro"

	updatedService, err := dao.UpdateService(UpdateServiceEntity{
		Id:          &serviceId,
		Servicename: &servicename,
	})
	if err != nil || updatedService.Id != serviceId || updatedService.Servicename != servicename {
		t.Fatalf("no item(updated service error) %#v", err)
	}

	getService1, err := dao.GetService(serviceId)
	if err != nil || getService1 == nil {
		t.Fatalf("no item(create service error) %#v", err)
	}
	if diff := cmp.Diff(*getService1, createdService); diff == "" {
		t.Fatalf("failed test(created service is wrong) %#v", err)
	}
	if getService1.Id != serviceId || getService1.Servicename != servicename {
		t.Fatalf("no item(updated service error) %#v", err)
	}

}

func TestPOST_DELETE_GET(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), dynamoLocalEndpoint)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	id, _ := uuid.NewUUID()
	serviceId := id.String()
	getService0, err := dao.GetService(serviceId)
	if err != nil || getService0 != nil {
		t.Fatalf("failed test %#v", err)
	}

	if _, err := dao.DeleteService(serviceId); err != nil {
		t.Fatalf("failed test(need to initialize dynamodb local) %#v", err)
	}
	getService1, err := dao.GetService(serviceId)
	if err != nil || getService1 != nil {
		t.Fatalf("no item(create service error) %#v", err)
	}

}
