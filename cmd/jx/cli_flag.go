// +build !kong

package main

import (
	"flag"
	"os"

	"foxygo.at/jsonnext"
)

func parseCLI() *config {
	c := &config{}
	c.Config = *jsonnext.ConfigFlags(flag.CommandLine)
	flag.Parse()
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	}
	c.Filename = flag.Args()[0]
	return c
}
