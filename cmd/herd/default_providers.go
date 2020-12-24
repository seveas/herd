package main

// These imports are explicitely ordered/grouped to make sure we register
// providers in the desired order
import (
	_ "github.com/seveas/herd/provider/json"
	_ "github.com/seveas/herd/provider/plain"

	_ "github.com/seveas/herd/provider/aws"
	_ "github.com/seveas/herd/provider/cache"
	_ "github.com/seveas/herd/provider/consul"
	_ "github.com/seveas/herd/provider/http"
	_ "github.com/seveas/herd/provider/prometheus"
)
