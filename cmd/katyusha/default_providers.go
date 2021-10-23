package main

// These imports are explicitely ordered/grouped to make sure we register
// providers in the desired order
import (
	// The basics for introperability with openssh and putty
	_ "github.com/seveas/katyusha/provider/known_hosts"
	_ "github.com/seveas/katyusha/provider/putty"

	// Simple file based providers
	_ "github.com/seveas/katyusha/provider/json"
	_ "github.com/seveas/katyusha/provider/plain"

	// Network based ones
	_ "github.com/seveas/katyusha/provider/cache"
	_ "github.com/seveas/katyusha/provider/consul"
	_ "github.com/seveas/katyusha/provider/http"
	_ "github.com/seveas/katyusha/provider/prometheus"

	// Cloud providers
	_ "github.com/seveas/katyusha/provider/aws"
	_ "github.com/seveas/katyusha/provider/azure"
	_ "github.com/seveas/katyusha/provider/google"
	_ "github.com/seveas/katyusha/provider/transip"

	// The sky is the limit!
	_ "github.com/seveas/katyusha/provider/plugin"
)
