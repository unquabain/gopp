package main

import (
	"github.com/apex/log"
	"github.com/spf13/pflag"
)

type Options struct {
	debug bool
}

func (o *Options) Parse() {
	pflag.BoolVarP(&o.debug, "debug", `g`, false, "enable debug mode")
	pflag.Parse()
	if o.debug {
		log.SetLevel(log.DebugLevel)
	}
}

func (o *Options) Inputs() *Inputs {
	return NewFilenames(pflag.Args())
}
