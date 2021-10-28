//go:build never
// +build never

package main

import (
	"fmt"
	"github.com/seveas/herd"
)

func main() {
	fmt.Println(herd.Version())
}
