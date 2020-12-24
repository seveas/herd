package main

// These imports are explicitely ordered/grouped to make sure we register
// providers in the desired order
import (
	_ "github.com/seveas/katyusha/provider/known_hosts"

	_ "github.com/seveas/katyusha/provider/json"
	_ "github.com/seveas/katyusha/provider/plain"

	_ "github.com/seveas/katyusha/provider/aws"
	_ "github.com/seveas/katyusha/provider/cache"
	_ "github.com/seveas/katyusha/provider/consul"
	_ "github.com/seveas/katyusha/provider/http"
	_ "github.com/seveas/katyusha/provider/prometheus"
)
