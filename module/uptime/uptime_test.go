package uptime

import (
	"time"
	"testing"

	"github.com/innectic/starlink-outages/module"

	pb "github.com/starlink-community/starlink-grpc-go/pkg/spacex.com/api/device"
)

var messageChan = make(chan module.ModuleMessage)
var first = true
var satellites = false

type RPCHandler struct {
}

func (rpc RPCHandler) GetStatus() (*pb.DishGetStatusResponse, error) {
	if first {
		first = false
		if satellites {
			return &pb.DishGetStatusResponse{
				SecondsToFirstNonemptySlot: 5,
			}, nil
		}

		return &pb.DishGetStatusResponse{
			PopPingDropRate: 1,
		}, nil
	}

	return &pb.DishGetStatusResponse{
		PopPingDropRate: 0,
	}, nil
}


func TestRunForBetaDownTime(t *testing.T) {
	satellites = false
	first = true

	go func() {
		<-messageChan
	}()

	last := LastData{
		PopPingDropRate: 0,
		OutageStart: time.Now(),
		OutageEnd: time.Now(),
		Cause: DishyOnline,
	}
	r := RPCHandler{ }

	module := UptimeModule{
		c: messageChan,
		r: r,
		lastData: &last,
	}

	if err := module.Run(); err != nil {
		t.Fatalf("faliled to run uptime: %v", err)
		return
	}

	time.Sleep(2 * time.Second)

	if err := module.Run(); err != nil {
		t.Fatalf("faliled to run uptime: %v", err)
	}
}

func TestRunForNoSatellites(t *testing.T) {
	satellites = true
	first = true

	go func() {
		<-messageChan
	}()

	t.Log("Starting satellite")
	last := LastData{
		PopPingDropRate: 0,
		OutageStart: time.Now(),
		OutageEnd: time.Now(),
		Cause: DishyOnline,
	}
	r := RPCHandler{ }

	module := UptimeModule{
		c: messageChan,
		r: r,
		lastData: &last,
	}

	if err := module.Run(); err != nil {
		t.Fatalf("faliled to run uptime: %v", err)
		return
	}

	time.Sleep(2 * time.Second)

	if err := module.Run(); err != nil {
		t.Fatalf("faliled to run uptime: %v", err)
	}
}
