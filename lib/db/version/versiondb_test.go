package versiondb

import (
	"os"
	"testing"
)

var dynamoLocalEndpoint string = "http://localhost:8027"

func TestGetServiceSuccess(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("VERSIONTABLENAME", "swagger-dev-swagger-dynamo-versioninfo")

	dao, err := NewDaoWithRegionAndEndpoint(os.Getenv("VERSIONTABLENAME"), os.Getenv("AWS_DEFAULT_REGION"), dynamoLocalEndpoint)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	serviceId := "66a36e77-fd00-3779-8097-17841f998f4d"

	versions, err := dao.GetAllVersions(serviceId)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	// fmt.Printf("%+v\n", versions)

	for _, v := range versions {
		if v.ID != serviceId {
			t.Fatalf("failed test %#v", err)
		}
	}

}
