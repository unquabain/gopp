package main

import (
	"bytes"
	"io"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/apex/log"
	"gopkg.in/yaml.v3"
)

type RTemplate struct {
	template *template.Template
	context  any
}

func NewRTemplate(input io.Reader, context any) (*RTemplate, error) {
	if context == nil {
		context = make(map[string]any)
	}
	r := &RTemplate{
		context: context,
	}
	if err := r.open(input); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RTemplate) FuncMap() template.FuncMap {
	fm := sprig.TxtFuncMap()
	fm["include"] = r.Include
	return fm
}

func (r *RTemplate) Include(fname string, context any) (string, error) {
	logger := log.WithField("fname", fname).WithField("context", context)
	logger.Debug("include")
	if context == nil {
		context = r.context
	}
	input, err := os.Open(fname)
	if err != nil {
		logger.WithError(err).Error("failed to open file")
		return "", err
	}
	defer input.Close()
	rt, err := NewRTemplate(input, context)
	if err != nil {
		logger.WithError(err).Error("failed to create template")
		return "", err
	}
	return rt.Render()
}

func (r *RTemplate) open(f io.Reader) error {
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if string(b[:4]) == "---\n" {
		parts := bytes.SplitN(b, []byte("\n---\n"), 2)
		if len(parts) == 2 {
			if err := yaml.Unmarshal(parts[0], r.context); err != nil {
				return err
			}
			b = parts[1]
		}
	}
	r.template, err = template.New("").Funcs(r.FuncMap()).Parse(string(b))
	return err
}

func (r *RTemplate) Render() (string, error) {
	var b bytes.Buffer
	if err := r.template.Execute(&b, r.context); err != nil {
		return "", err
	}
	return b.String(), nil
}
