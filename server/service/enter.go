package service

import (
	"github.com/slyrx/gin_exam_system/server/service/example"
	"github.com/slyrx/gin_exam_system/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
