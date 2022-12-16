package main

import (
	// Import the provider you wish to serve over grpc
	_ "github.com/seveas/herd/provider/plugin/testdata/provider/ci_dataloader"

	// And the helper library to serve it
	"github.com/seveas/herd/provider/plugin/server"
)

func main() {
	if err := server.ProviderPluginServer("ci_dataloader"); err != nil {
		panic(err)
	}
}
