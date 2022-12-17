package main

import (
	// Import the provider you wish to serve over grpc, and the helper library to serve it
	_ "github.com/seveas/herd/provider/azure"
	"github.com/seveas/herd/provider/plugin/server"
)

func main() {
	if err := server.ProviderPluginServer("azure"); err != nil {
		panic(err)
	}
}
