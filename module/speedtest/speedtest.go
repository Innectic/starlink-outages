package speedtest

import (
	"github.com/innectic/starlinkoutages/log"
	"github.com/innectic/starlinkoutages/rpc"
	"github.com/innectic/starlinkoutages/module"

	"github.com/showwin/speedtest-go/speedtest"
)

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
	user, _ := speedtest.FetchUserInfo()

	serverList, _ := speedtest.FetchServerList(user)
	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)

		
	}

	return nil, nil
}
