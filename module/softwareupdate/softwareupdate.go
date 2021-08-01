package softwareupdate

import (
	"fmt"
	"time"

	"github.com/innectic/starlinkoutages/log"
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
	var lastVersion string = ""
	if last != nil {
		lastVersion = last.(string)
	}

	dishy, err := m.r.GetStatus()
	if err != nil {
		return lastVersion, err
	}

	dishyVersion := dishy.DeviceInfo.SoftwareVersion

	// If we don't currently know what version is running, skip this and wait until the next one.
	if lastVersion == "" {
		return dishyVersion, nil
	}

	// Since we do have a version, then compare previous to the current one
	if lastVersion != dishyVersion {
		// Version updated.
		log.Info(fmt.Sprintf("Dishy version updated! Old: %s - New: %s", lastVersion, dishyVersion))
		m.c <- module.ModuleMessage{
			Message: GetMessage(time.Now().Format("01/02/2006 15:04:05 MST"), lastVersion, dishyVersion),
		}
		return dishyVersion, nil
	}

	return lastVersion, nil
}
