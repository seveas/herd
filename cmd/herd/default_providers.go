package main

// These imports are explicitly ordered/grouped to make sure we register
// providers in the desired order
import (
	// The basics for introperability with openssh and putty
	_ "github.com/seveas/herd/provider/known_hosts"
	_ "github.com/seveas/herd/provider/putty"

	// Simple file based providers
	_ "github.com/seveas/herd/provider/json"
	_ "github.com/seveas/herd/provider/plain"

	// Network based ones
	_ "github.com/seveas/herd/provider/cache"
	_ "github.com/seveas/herd/provider/consul"
	_ "github.com/seveas/herd/provider/http"
	_ "github.com/seveas/herd/provider/prometheus"

	// Disabled for now, as it doesn't build on windows
	// _ "github.com/seveas/herd/provider/tailscale"

	// Cloud providers
	_ "github.com/seveas/herd/provider/aws"
	_ "github.com/seveas/herd/provider/azure"
	_ "github.com/seveas/herd/provider/google"
	_ "github.com/seveas/herd/provider/transip"

	// The sky is the limit!
	_ "github.com/seveas/herd/provider/plugin"
)
