package rpc

import (
	"time"
	"context"
	pb "github.com/starlink-community/starlink-grpc-go/pkg/spacex.com/api/device"
	"google.golang.org/grpc"
)

const rpcTimeout = 5 * time.Second

type RPCHandler interface {
	GetStatus() (*pb.DishGetStatusResponse, error)
}

type RPCHandlerImpl struct {
	Address string
}

func NewRPCHandler(address string) *RPCHandlerImpl {
	return &RPCHandlerImpl{
		Address: address,
	}
}

func (rpc RPCHandlerImpl) GetStatus() (*pb.DishGetStatusResponse, error) {
	conn, err := grpc.Dial(rpc.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := pb.NewDeviceClient(conn)
	in := new(pb.Request)

	in.Request = &pb.Request_GetStatus{}
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()

	r, err := c.Handle(ctx, in)
	if err != nil {
		return nil, err
	}

	return r.GetDishGetStatus(), nil
}
