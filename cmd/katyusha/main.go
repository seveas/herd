package main

import (
	"fmt"

	"github.com/seveas/katyusha"
)

func main() {
	providers := katyusha.LoadProviders()
	hosts := providers.GetHosts("ops-shell-*", katyusha.HostAttributes{})
	fmt.Printf("%d hosts\n", len(hosts))
	fmt.Printf("%s\n", hosts)

	runner := katyusha.NewRunner(hosts)
	hi := runner.Run("id; uptime; hostname")
	fmt.Printf("%v\n", hi.Results)
}
