package uptime

import (
	"fmt"
	"time"

	"github.com/innectic/starlinkoutages/log"
	"github.com/innectic/starlinkoutages/rpc"
	"github.com/innectic/starlinkoutages/module"

	pb "github.com/starlink-community/starlink-grpc-go/pkg/spacex.com/api/device"
)

var crap int = 0

type DishyStatus string
var (
	DishyOnline DishyStatus = "ONLINE"
	DishyBetaDownTime DishyStatus = "BETA DOWNTIME / OTHER OUTAGES"
	DishyNoSatellites DishyStatus = "NO SATELLITES"
)

type UptimeModule struct {
	c chan module.ModuleMessage
	r rpc.RPCHandler
}

func NewUptimeModule(c chan module.ModuleMessage, r rpc.RPCHandler) UptimeModule {
	return UptimeModule{
		c: c,
		r: r,
	}
}

func (u UptimeModule) Init() (module.ModuleDefinition, error) {
	def := module.ModuleDefinition{
		Name: "Uptime Monitor",
		Description: "Monitors uptime and outage duration / causes",
		Frequency: 1 * time.Second,
	}

	return def, nil
}

func determineStatus(dishy pb.DishGetStatusResponse) DishyStatus {
	if dishy.PopPingDropRate == 1 {
		return DishyBetaDownTime
	} else if dishy.SecondsToFirstNonemptySlot > 1 {
		return DishyNoSatellites
	}
	return DishyOnline
}

func (u UptimeModule) Run(last interface{}) (interface{}, error) {
	var l LastData
	if last == nil {
		l = defaultLastData()
	} else {
		l = last.(LastData)
	}

	dishy, err := u.r.GetStatus()
	if err != nil {
		return l, err
	}

	currentState := determineStatus(*dishy)
	if crap == 30 {
		log.Info("Current State: " + currentState)
		crap = 0
	}
	crap++

	// Check if starlink is currently online
	if currentState == DishyOnline {
		// Starlink is currently online. Did we just finish an outage?

		// TODO: Support outage chaining
		if l.Cause == DishyOnline {
			// Last state was also connected, therefor we are perfectly fine.
			return l, nil
		}

		// Collect data to be published
		l.Ended()
		duration := l.Duration()
		friendlyStart := l.FriendlyStartTime()

		log.Info("OUTAGE COMPLETE! Duration: " + fmt.Sprint(duration))

		// Last state was not connected, so we just finished an outage.
		u.c <- module.ModuleMessage{
			Message: GetMessage(friendlyStart, duration, l.Cause, 0, 0, nil),
		}
	} else {
		// Dishy is not currently online. Start collecting data.
		if !l.Collecting {
			l.Started()
			log.Info(fmt.Sprintf("OUTAGE STARTED! Cause: %s - Start time: %s", currentState, l.FriendlyStartTime()))
		}
	}
	l.Cause = currentState

	return l, nil
}
