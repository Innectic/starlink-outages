package softwareupdate

import (
	"time"

	"github.com/innectic/starlinkoutages/rpc"
	"github.com/innectic/starlinkoutages/module"
)

type SoftwareUpdateModule struct {
	c chan module.ModuleMessage
	r rpc.RPCHandler
}

func NewSoftwareUpdateModule(c chan module.ModuleMessage, r rpc.RPCHandler) SoftwareUpdateModule {
	return SoftwareUpdateModule{
		c: c,
		r: r,
	}
}

func (m SoftwareUpdateModule) Init() (module.ModuleDefinition, error) {
	def := module.ModuleDefinition{
		Name: "Software Update Monitor",
		Description: "Monitors the software version of Dishy for updates",
		Frequency: 1 * time.Minute,
	}

	return def, nil
}

func (m SoftwareUpdateModule) Run(last interface{}) (interface{}, error) {
	return nil, nil
}
