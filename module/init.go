package module

import (
	"time"
)

type ModuleMessage struct {
	Message string
}

type ModuleDefinition struct {
	Name string
	Description string
	Frequency time.Duration
}

type Module interface {
	Init() (ModuleDefinition, error)
	Run(last interface{}) (interface{}, error)
}
