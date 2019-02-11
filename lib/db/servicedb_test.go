package servicedb

import (
	// "fmt"
	"os"
	"testing"
)

func TestGetServiciesSuccess(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), "http://localhost:8000")
	// dao, err := NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"));
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	servicies, err := dao.GetServiceList()
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(servicies)
}

func TestGetServiceSuccess(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), "http://localhost:8000")
	// dao, err := NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"));
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	serviceId := "66a36e77-fd00-3779-8097-17841f998f4d"

	service, err := dao.GetService(serviceId)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if service.Id != serviceId {
		t.Fatalf("invalid id %v", service)
	}
}

func TestGetServiceShouldReturnNil(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), "http://localhost:8000")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	serviceId := "66a36e"

	service, err := dao.GetService(serviceId)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if service != nil {
		t.Fatalf("invalid id %v", service)
	}
}

func TestUpdateServiceSuccess(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), "http://localhost:8000")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	servicename := "test"
	var lastupdated int64 = 111111
	id := "66a36e77-fd00-3779-8097-17841f998f4d"

	result, err := dao.UpdateService(UpdateServiceEntity{
		Id:          &id,
		Servicename: &servicename,
		Lastupdated: &lastupdated,
	})
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Logf("%+v\n", result)
	if result.Id != id {
		t.Fatalf("incorrect id")
	}
	if result.Servicename != servicename {
		t.Fatalf("incorrect servicename")
	}
	if result.Lastupdated != lastupdated {
		t.Fatalf("incorrect lastupdated")
	}
}

func TestUpdateServiceNoUpdateShouldReturnError(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), "http://localhost:8000")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	id := "66a36e77-fd00-3779-8097-17841f998f4d"

	result, err := dao.UpdateService(UpdateServiceEntity{
		Id: &id,
	})
	if err == nil {
		t.Fatalf("should return error %#v", result)
	}
}

func TestUpdateServiceIncorrectIdShouldReturnError(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), "http://localhost:8000")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	id := "1"
	servicename := "test"

	result, err := dao.UpdateService(UpdateServiceEntity{
		Id:          &id,
		Servicename: &servicename,
	})
	if err == nil {
		t.Fatalf("should return error %#v", result)
	}
	t.Logf("failed test %#v", err)
}

func TestPostServiceSuccess(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), "http://localhost:8000")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	result, err := dao.CreateService(ServiceEntity{
		Id:            "createid",
		Servicename:   "testservice",
		Latestversion: "1.2.3",
		Lastupdated:   53,
	})

	if err != nil {
		t.Fatalf("failed test %#v", err.Error())
	}

	if result.Id != "" {
		t.Fatal("return value should be nothing")
	}

}

func TestDeleteServiceSuccess(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("SERVICETABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), "http://localhost:8000")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	id := "createid"
	result, err := dao.DeleteService(id)

	if err != nil {
		t.Fatalf("failed test %#v", err.Error())
	}
	if result.Id != id {
		t.Fatalf("incorrect id")
	}
}
