//go:build never
// +build never

package main

import (
	"fmt"
	"github.com/seveas/katyusha"
)

func main() {
	fmt.Println(katyusha.Version())
}
