package main

import (
	// Import the provider you wish to serve over grpc, and the helper library to serve it
	_ "github.com/seveas/herd/provider/aws"
	"github.com/seveas/herd/provider/plugin/server"
)

func main() {
	if err := server.ProviderPluginServer("aws"); err != nil {
		panic(err)
	}
}
