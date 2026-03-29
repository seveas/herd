package main

import (
	"os"

	// Import the provider you wish to serve over grpc
	_ "github.com/seveas/herd/provider/plugin/testdata/provider/ci"

	// And the helper library to serve it
	"github.com/seveas/herd/provider/plugin/server"
)

func main() {
	if err := server.ProviderPluginServer("ci", os.Args[1]); err != nil {
		panic(err)
	}
}
