package speedtest

import (
	"fmt"
	"time"
	"github.com/innectic/starlinkoutages/log"
	"github.com/innectic/starlinkoutages/rpc"
	"github.com/innectic/starlinkoutages/module"

	"github.com/showwin/speedtest-go/speedtest"
)

var resetLast = false

type SpeedtestModule struct {
	r rpc.RPCHandler
	c chan module.ModuleMessage
}

func (m SpeedtestModule) Init() (module.ModuleDefinition, error) {
	def := module.ModuleDefinition{
		Name: "Speedtest",
		Description: "Periodically run speedtests to track speed stability",
		Frequency: 1 * time.Hour,
	}

	return def, nil
}

func (m SpeedtestModule) Run(last interface{}) (interface{}, error) {
	var l LastData
	if last == nil {
		l = defaultLastData()
	} else {
		l = last.(LastData)
	}

	if resetLast == true {
		resetLast = false
		l.Reset()
	}

	user, _ := speedtest.FetchUserInfo()

	serverList, _ := speedtest.FetchServerList(user)
	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)

		log.Info(fmt.Sprintf("Latency: %s, download: %f, upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed))
		l.Result(int(s.Latency), int(s.ULSpeed), int(s.DLSpeed))
	
		m.c <- module.ModuleMessage{
			Message: EachHour(int(s.Latency), int(s.ULSpeed), int(s.DLSpeed)),
		}
	}

	return l, nil
}

func (m SpeedtestModule) Reset() {
	resetLast = true

	// Tweet daily results
	// TODO
}
