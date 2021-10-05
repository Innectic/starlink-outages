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

func NewSpeedtestModule(c chan module.ModuleMessage, r rpc.RPCHandler) SpeedtestModule {
	return SpeedtestModule{
		c: c,
		r: r,
	}
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

		downloadAvg, downloadLow, downloadHigh := l.Download()
		uploadAvg, uploadLow, uploadHigh := l.Upload()
		latencyAvg, latencyLow, latencyHigh := l.Latency()
		failedTests := l.Failed()

		message := Daily(downloadAvg, downloadLow, downloadHigh, uploadAvg, uploadLow, uploadHigh, latencyAvg, latencyLow, latencyHigh, l.TotalRuns, failedTests)
		m.c <- module.ModuleMessage{
			Message: message,
		}

		l.Reset()
	}
	l.TotalRuns += 1

	user, err := speedtest.FetchUserInfo()
	if err != nil {
		return l, nil
	}

	serverList, err := speedtest.FetchServerList(user)
	if err != nil {
		return l, nil
	}
	targets, err := serverList.FindServer([]int{})
	if err != nil {
		return l, nil
	}

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)

		log.Info(fmt.Sprintf("Latency: %s, download: %f, upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed))
		l.Result(int(s.Latency), int(s.ULSpeed), int(s.DLSpeed))
		l.SuccessfulRuns += 1

		m.c <- module.ModuleMessage{
			Message: EachHour(int(s.Latency), int(s.ULSpeed), int(s.DLSpeed)),
		}
	}

	return l, nil
}

func (m SpeedtestModule) Reset() {
	resetLast = true
}
