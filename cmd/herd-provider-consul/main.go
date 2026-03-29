package main

import (
	"os"

	// Import the provider you wish to serve over grpc, and the helper library to serve it
	_ "github.com/seveas/herd/provider/consul"
	"github.com/seveas/herd/provider/plugin/server"
)

func main() {
	if err := server.ProviderPluginServer("consul", os.Args[1]); err != nil {
		panic(err)
	}
}
