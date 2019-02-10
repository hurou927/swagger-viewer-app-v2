package servicedb

import (
    // "fmt"
    "testing"
    "os"
)

func TestGetServicies (t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION","ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")
	
	dao, err := NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"));
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	servicies, err := dao.serviceRepositoryDao.GetServiceList();
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	t.Log(servicies);

}