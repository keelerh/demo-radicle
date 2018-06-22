package main

import (
	"flag"
	"log"

	pb "github.com/keelerh/radicle-demo/protos"
	"github.com/keelerh/radicle-demo/pkg/radicle"
	"github.com/keelerh/radicle-demo/pkg/radicle/common"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/keelerh/radicle-demo/pkg/helloworld"
)

var (
	fGRPCServerPort = common.FlagSet().Int(
		"grpc_server_port",
		8081,
		"Port on which the gRPC server should listen",
	)
	fHTTPServerPort = common.FlagSet().Int(
		"http_server_port",
		8080,
		"Port on which the HTTP server should listen",
	)
)

func main() {
	flag.Parse()

	interceptors := grpc_middleware.WithUnaryServerChain(
		grpc_prometheus.UnaryServerInterceptor,
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_recovery.UnaryServerInterceptor(),
	)
	srv, err := radicle.Listen(
		radicle.Grpc(*fGRPCServerPort),
		radicle.Http(*fHTTPServerPort),
		radicle.GrpcServerOption(interceptors),
	)
	if err != nil {
		log.Fatalf("failed to start listeners: %v", err)
	}

	svc := &helloworld.GreeterService{}
	pb.RegisterGreeterServer(srv.Grpc(), svc)

	if err := srv.Serve(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
