package main

import (
	// Import the provider you wish to serve over grpc
	_ "github.com/seveas/katyusha/provider/plugin/testdata/provider/ci"

	// And the helper library to serve it
	"github.com/seveas/katyusha/provider/plugin/server"
)

func main() {
	if err := server.ProviderPluginServer("ci"); err != nil {
		panic(err)
	}
}
