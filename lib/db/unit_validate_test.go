package servicedb

import (
	"testing"
)

func TestServiceNameValidate(t *testing.T) {
	if ValidateServiceName("auth") != true {
		t.Fatalf("failed test")
	}
}

func TestServiceNameValidate2(t *testing.T) {
	if ValidateServiceName("auth audit") != false {
		t.Fatalf("failed test")
	}
}
