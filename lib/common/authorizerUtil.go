package common

import "gopkg.in/yaml.v2"

type AuthorizerConfig struct {
	WhitelistIP []string `json:"whitelist_ip" yaml:"whitelist_ip" validate:"required"`
	BlacklistIP []string `json:"blacklist_ip" yaml:"blacklist_ip" validate:"required"`
}

func ParseAuthorizerConfig(config string) (AuthorizerConfig, error) {
	var authorizerConfig AuthorizerConfig
	if err := yaml.Unmarshal([]byte(config), &authorizerConfig); err != nil {
		return AuthorizerConfig{}, err
	}
	return authorizerConfig, nil
}
