package main

import (
	"context"
	"fmt"
	"gitlab.com/firestart/ignition/examples/application/rpc/pb"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/x/application/extensions/grpc"
	"gitlab.com/firestart/ignition/x/application/presets"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	app := presets.NewRpcApp("test")

	// Initialize the application
	grpc.MustAddService(app, pb.Greeter_ServiceDesc, MustProvideGreeterService(app.Injector))

	// Start the application
	app.Run()
}

type GreeterService struct {
}

func NewGreeterService() GreeterService {
	return GreeterService{}
}

func ProvideGreeterService(inj *injector.Injector) (GreeterService, error) {
	return NewGreeterService(), nil
}

func MustProvideGreeterService(inj *injector.Injector) GreeterService {
	service, err := ProvideGreeterService(inj)
	if err != nil {
		panic(err)
	}
	return service
}

func (g GreeterService) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "Hello " + request.Name,
	}, nil
}

func (g GreeterService) Panic(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	switch request.Name {
	case "Error":
		return nil, fmt.Errorf("error")
	case "Code":
		return nil, status.Error(codes.NotFound, "not found")
	default:
		panic(request.Name)
	}
}

var _ pb.GreeterServer = GreeterService{}
