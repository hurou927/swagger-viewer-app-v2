package common

import (
	"fmt"
	"testing"
)

func TestAuthorizerSuccess(t *testing.T) {
	configyaml := `
whitelist_ip:
  - 192.168.11.1
`
	conf, err := ParseAuthorizerConfig(configyaml)
	fmt.Println(conf, err)
}
