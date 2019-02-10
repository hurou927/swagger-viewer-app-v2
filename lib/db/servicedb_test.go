package servicedb

import (
    "fmt"
    "testing"
    "os"
)

func TestGetServicies (t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION","ap-northeast-1")
	os.Setenv("SERVICETABLENAME", "swagger-dev-swagger-dynamo-serviceinfo")
	
	serviceDao, err := NewDaoDefaultConfig(os.Getenv("SERVICETABLENAME"));
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	servicies, err := serviceDao.getServiceList();
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	fmt.Println(servicies);

}