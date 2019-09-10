package test

import (
	"flag"
	"github.com/stainour/test10/api"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:9765", "The server address in the format of host:port")
)

func GetClient() (api.UrlShortenerClient, error) {
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return api.NewUrlShortenerClient(conn), nil
}
