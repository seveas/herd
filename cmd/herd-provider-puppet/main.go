package main

import (
	"os"

	// Import the provider you wish to serve over grpc, and the helper library to serve it
	"github.com/seveas/herd/provider/plugin/server"
	_ "github.com/seveas/herd/provider/puppet"
)

func main() {
	if err := server.ProviderPluginServer("puppet", os.Args[1]); err != nil {
		panic(err)
	}
}
