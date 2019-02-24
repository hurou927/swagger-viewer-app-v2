package common

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v2"
)

// Format represents file-format (Yaml/JSON)
type Format int

const (
	Yml Format = iota
	Json
)

type Swagger struct {
	Swagger string `json:"swagger" validate:"required"`
	Info    struct {
		Version string `json:"version" validate:"required"`
	} `json:"info" validate:"required"`
}

func ValidateSwagger(format Format, contents string) (Swagger, error) {
	var swagger Swagger
	if format == Yml {
		if err := yaml.Unmarshal([]byte(contents), &swagger); err != nil {
			return Swagger{}, NewError(20001, "Swagger(YML) Unmarshal Error", err)
		}
	} else {
		if err := json.Unmarshal([]byte(contents), &swagger); err != nil {
			return Swagger{}, NewError(20001, "Swagger(JSON) Unmarshal Error", err)
		}
	}

	versions := strings.Split(strings.TrimSpace(swagger.Info.Version), ".")

	if l := len(versions); l >= 1 && l <= 3 {
		for i := 0; i < l; i++ {
			versions[i] = strings.TrimSpace(versions[i])
		}
		swagger.Info.Version = strings.Join(versions, ".")
		return swagger, nil
	}

	return Swagger{}, NewError(20002, "Version Format Error", nil)

}
