package main

import (
	// Import the provider you wish to serve over grpc, and the helper library to serve it
	"github.com/seveas/herd/provider/plugin/server"
	_ "github.com/seveas/herd/provider/prometheus"
)

func main() {
	if err := server.ProviderPluginServer("prometheus"); err != nil {
		panic(err)
	}
}
