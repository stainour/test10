package test

import (
	"fmt"
	"github.com/stainour/test10/api"
	"github.com/stainour/test10/bootstrap"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = 9765

func StartServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service, err := bootstrap.InitializeService()

	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	grpcServer := grpc.NewServer()

	api.RegisterUrlShortenerServer(grpcServer, service)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
