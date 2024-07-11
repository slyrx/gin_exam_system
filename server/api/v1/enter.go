package v1

import (
	"github.com/slyrx/gin_exam_system/server/api/v1/example"
	"github.com/slyrx/gin_exam_system/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
