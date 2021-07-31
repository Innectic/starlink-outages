package main

import (
	"fmt"
	"time"
	"github.com/innectic/starlinkoutages/rpc"
	"github.com/innectic/starlinkoutages/log"
	"github.com/innectic/starlinkoutages/module"
	"github.com/innectic/starlinkoutages/module/uptime"
)

const (
	starlinkAddr = "192.168.100.1:9200"
)

func main() {
	h := rpc.NewRPCHandler(starlinkAddr)

	log.Info("Loading modules...")

	c := make(chan module.ModuleMessage)

	up := uptime.NewUptimeModule(c, *h, &uptime.LastData{})
	modules := []module.Module{ up }

	for _, mod := range modules {
		go func() {
			def, _ := mod.Init()

			for {
				mod.Run()
				time.Sleep(def.Frequency)
			}
		}()
	}

	// Start watching the message queue
	for {
		m := <-c
		fmt.Println(m)
	}
}
