package main

import (
	"fmt"

	"github.com/seveas/herd"
)

func main() {
	providers := herd.LoadProviders()
	hosts := providers.GetHosts("ops-shell-*", herd.HostAttributes{})
	fmt.Printf("%d hosts\n", len(hosts))
	fmt.Printf("%s\n", hosts)

	runner := herd.NewRunner(hosts)
	hi := runner.Run("id; uptime; hostname")
	fmt.Printf("%v\n", hi.Results)
}
