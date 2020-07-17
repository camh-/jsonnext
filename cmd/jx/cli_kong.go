// +build kong

package main

import "github.com/alecthomas/kong"

func parseCLI() *config {
	c := &config{}
	kong.Parse(c)
	return c
}
