package main

import (
	"os"

	// Import the provider you wish to serve over grpc, and the helper library to serve it
	_ "github.com/seveas/herd/provider/example"
	"github.com/seveas/herd/provider/plugin/server"
)

func main() {
	if err := server.ProviderPluginServer("example", os.Args[1]); err != nil {
		panic(err)
	}
}
