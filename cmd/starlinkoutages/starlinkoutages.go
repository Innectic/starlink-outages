package main

import (
	"fmt"
	"time"
	"github.com/innectic/starlinkoutages/rpc"
	"github.com/innectic/starlinkoutages/log"
	"github.com/innectic/starlinkoutages/module"

	"github.com/innectic/starlinkoutages/module/uptime"
	"github.com/innectic/starlinkoutages/module/softwareupdate"
)

const (
	starlinkAddr = "192.168.100.1:9200"
)

func main() {
	h := rpc.NewRPCHandler(starlinkAddr)

	log.Info("Loading modules...")

	c := make(chan module.ModuleMessage)

	up := uptime.NewUptimeModule(c, *h)
	update := softwareupdate.NewSoftwareUpdateModule(c, *h)
	modules := []module.Module{ up, update }

	for _, m := range modules {
		go func(mod module.Module) {
			def, _ := mod.Init()
			log.Info(fmt.Sprintf("Module load: %s", def.Name))

			var res interface{} = nil
			var err error
			for {
				res, err = mod.Run(res)
				if err != nil {
					log.Error(fmt.Sprint("Module %s encountered an error: %v", def.Name, err));
					continue
				}
				time.Sleep(def.Frequency)
			}
		}(m)
	}

	// Start watching the message queue
	for {
		m := <-c
		fmt.Println(m)
	}
}
