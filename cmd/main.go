package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/keelerh/radicle-demo/pkg/helloworld"
	pb "github.com/keelerh/radicle-demo/protos"
	"google.golang.org/grpc"
)

var (
	fGRPCServerPort = flag.Int(
		"grpc_server_port",
		8081,
		"Port on which gRPC server should listen",
	)
	fHTTPServerPort = flag.Int(
		"http_server_port",
		8080,
		"Port on which the HTTP server should listen",
	)
)

func runGRPCServer(lis net.Listener) {
	srv := grpc.NewServer()

	svc := &helloworld.GreeterService{}
	pb.RegisterGreeterServer(srv, svc)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runHTTPServer() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := pb.RegisterGreeterHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("localhost:%d", *fGRPCServerPort),
		opts,
	); err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", *fHTTPServerPort), mux)
}

func main() {
	flag.Parse()

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *fGRPCServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go runGRPCServer(grpcListener)

	if err := runHTTPServer(); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
