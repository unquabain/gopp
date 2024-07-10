package main

import (
	"fmt"
	"os"
)

type TopError int

const (
	ErrCouldNotCreateRTemplate TopError = iota
	ErrCouldNotRenderRTemplate
	ErrCouldNotLoopThroughInputs
)

func (ret TopError) checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(-1 - int(ret))
}

func main() {
	opts := new(Options)
	opts.Parse()
	inputs := opts.Inputs()
	for inputs.Next() {
		rt, err := NewRTemplate(inputs.Reader(), nil)
		ErrCouldNotCreateRTemplate.checkErr(err)
		s, err := rt.Render()
		ErrCouldNotRenderRTemplate.checkErr(err)
		fmt.Print(s)
	}
	ErrCouldNotLoopThroughInputs.checkErr(inputs.Err())
}
