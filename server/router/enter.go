package router

import (
	"github.com/slyrx/gin_exam_system/server/router/example"
	"github.com/slyrx/gin_exam_system/server/router/system"
)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
