package main

import (
	"fmt"
	"time"
	"github.com/innectic/starlinkoutages/rpc"
	"github.com/innectic/starlinkoutages/log"
	"github.com/innectic/starlinkoutages/tweet"
	"github.com/innectic/starlinkoutages/config"
	"github.com/innectic/starlinkoutages/module"

	"github.com/innectic/starlinkoutages/module/uptime"
	"github.com/innectic/starlinkoutages/module/speedtest"
	"github.com/innectic/starlinkoutages/module/softwareupdate"
)

const (
	starlinkAddr = "192.168.100.1:9200"
)

func secondsUntilMidnight() int {
	tomorrow := time.Now().AddDate(0, 0, 1)
	midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 1, 0, time.Local)
	return int(midnight.Sub(time.Now()).Seconds())
}

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Error(fmt.Sprintf("failed to load configuration: %e\n", err))
		return
	}

	h := rpc.NewRPCHandler(starlinkAddr)

	log.Info("Setting up Twitter API...")
	twitter, _ := tweet.NewTweetQueue(cfg.ConsumerKey, cfg.ConsumerSecret, cfg.AccessToken, cfg.AccessSecret)
	go twitter.HandleTweetQueue()

	log.Info("Loading modules...")

	c := make(chan module.ModuleMessage)

	up := uptime.NewUptimeModule(c, *h)
	update := softwareupdate.NewSoftwareUpdateModule(c, *h)
	speedtest := speedtest.NewSpeedtestModule(c, *h)
	modules := []module.Module{ up, update, speedtest }

	for _, m := range modules {
		go func(mod module.Module) {
			def, _ := mod.Init()
			log.Info(fmt.Sprintf("Module load: %s", def.Name))

			var res interface{} = nil
			var err error
			for {
				res, err = mod.Run(res)
				if err != nil {
					log.Error(fmt.Sprint("Module %s encountered an error: %v", def.Name, err))
					continue
				}
				time.Sleep(def.Frequency)
			}
		}(m)
	}

	// Setup reset loop
	go func() {
		for {
			nextReset := secondsUntilMidnight()
			log.Info(fmt.Sprintf("Waiting %d seconds until next reset...", nextReset))
			time.Sleep(time.Duration(nextReset) * time.Second)
			log.Info("Issuing reset...")

			for _, m := range modules {
				m.Reset()
			}
			log.Info("Reset complete!")
		}
	}()

	// Start watching the message queue
	for {
		m := <-c
		twitter.QueueTweet(m.Message)
		fmt.Println(m)
	}
}
