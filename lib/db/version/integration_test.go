package versiondb

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

// var dynamoLocalEndpoint string = "http://localhost:8027"

func TestCREATE_GET(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("VERSIONTABLENAME", "swagger-dev-swagger-dynamo-versioninfo")
	os.Setenv("SWAGGER_BUCKET_NAME", "swagger-repository-test")
	dao, err := NewDaoWithEndpoints(
		os.Getenv("VERSIONTABLENAME"),
		AwsEndpoint{
			Region:   os.Getenv("AWS_DEFAULT_REGION"),
			Endpoint: dynamoLocalEndpoint,
		},
		AwsEndpoint{
			Region:   os.Getenv("AWS_DEFAULT_REGION"),
			Endpoint: s3LocalEndpoint,
		},
	)

	id, _ := uuid.NewUUID()
	serviceId := id.String()
	versions, err := dao.GetAllVersions(serviceId)
	if err != nil || versions != nil {
		t.Fatalf("failed test %#v", err)
	}

	bucketName := "swagger-repository-test"
	keyName := "keyname"
	contents := "swagger"
	tag := "tag"
	version := "10.2.23"
	requestEntity := VersionEntity{
		ID:          serviceId,
		Version:     version,
		Path:        keyName,
		Lastupdated: time.Now().Unix() * 1000,
		Enable:      true,
		Tag:         tag,
	}
	if _, err := dao.UploadVersion(requestEntity, bucketName, keyName, contents); err != nil {
		t.Fatalf("upload error %#v", err)
	}

	versions, err = dao.GetAllVersions(serviceId)
	if err != nil || versions == nil {
		t.Fatalf("failed test %#v", err)
	}
	versionInfo := versions[0]

	if diff := cmp.Diff(versionInfo, requestEntity); diff != "" {
		t.Fatalf("failed test(created service is wrong) %#v", err)
	}

}

func TestCREATE_UPDATE_GET(t *testing.T) {
	os.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")
	os.Setenv("VERSIONTABLENAME", "swagger-dev-swagger-dynamo-versioninfo")
	os.Setenv("SWAGGER_BUCKET_NAME", "swagger-repository-test")
	dao, err := NewDaoWithEndpoints(
		os.Getenv("VERSIONTABLENAME"),
		AwsEndpoint{
			Region:   os.Getenv("AWS_DEFAULT_REGION"),
			Endpoint: dynamoLocalEndpoint,
		},
		AwsEndpoint{
			Region:   os.Getenv("AWS_DEFAULT_REGION"),
			Endpoint: s3LocalEndpoint,
		},
	)

	id, _ := uuid.NewUUID()
	serviceId := id.String()
	versions, err := dao.GetAllVersions(serviceId)
	if err != nil || versions != nil {
		t.Fatalf("failed test %#v", err)
	}

	bucketName := "swagger-repository-test"
	keyName := "keyname"
	contents := "swagger"
	tag := "tag"
	version := "10.2.23"
	requestEntity := VersionEntity{
		ID:          serviceId,
		Version:     version,
		Path:        keyName,
		Lastupdated: time.Now().Unix() * 1000,
		Enable:      true,
		Tag:         tag,
	}
	if _, err := dao.UploadVersion(requestEntity, bucketName, keyName, contents); err != nil {
		t.Fatalf("upload error %#v", err)
	}

	tag2 := "tag2"
	updatedRequestEntity := VersionEntity{
		ID:          serviceId,
		Version:     version,
		Path:        keyName,
		Lastupdated: time.Now().Unix() * 1000,
		Enable:      false,
		Tag:         tag2,
	}

	if _, err := dao.UpdateVersion(updatedRequestEntity); err != nil {
		t.Fatalf("update error %#v", err)
	}

	versions, err = dao.GetAllVersions(serviceId)
	if err != nil || versions == nil {
		t.Fatalf("failed test %#v", err)
	}
	versionInfo := versions[0]

	if diff := cmp.Diff(versionInfo, requestEntity); diff == "" {
		t.Fatalf("failed test(created service is wrong) %#v", err)
	}
	if diff := cmp.Diff(versionInfo, updatedRequestEntity); diff != "" {
		t.Fatalf("failed test(created service is wrong) %#v", err)
	}

}
