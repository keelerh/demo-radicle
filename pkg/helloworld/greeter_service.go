package helloworld

import (
	"context"
	"fmt"
	"log"

	pb "github.com/keelerh/radicle-demo/protos"
)

type GreeterService struct{}

func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	//helloCount.WithLabelValues(req.Name).Add(1)
	log.Printf("%s: %s: %s", "grpc", "SayHello", req.Name)
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}
