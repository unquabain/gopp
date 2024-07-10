package main

import (
	"io"
	"os"

	"github.com/apex/log"
)

type Inputs struct {
	filenames []string
	next      io.ReadCloser
	error     error
}

func NewFilenames(args []string) *Inputs {
	if len(args) == 0 {
		return &Inputs{
			filenames: []string{`-`},
		}
	}
	return &Inputs{
		filenames: args,
	}
}

func (i *Inputs) Next() bool {
	logger := log.WithField(`filenames`, i.filenames)
	logger.Debug(`next`)
	if i.error != nil {
		logger.WithError(i.error).Error(`error`)
		return false
	}
	if i.next != nil {
		if i.error = i.next.Close(); i.error != nil {
			logger.WithError(i.error).Error(`failed to close`)
			return false
		}
		i.next = nil
	}
	if len(i.filenames) == 0 {
		return false
	}
	var fname string
	fname, i.filenames = i.filenames[0], i.filenames[1:]
	if fname == `-` {
		i.next = os.Stdin
		return true
	}
	i.next, i.error = os.Open(fname)
	if err := i.error; err != nil {
		logger.WithError(err).Error(`failed to open`)
	}
	return i.error == nil
}

func (i *Inputs) Err() error {
	return i.error
}

func (i *Inputs) Reader() io.Reader {
	return i.next
}
