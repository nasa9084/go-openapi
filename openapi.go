package openapi

import (
	"bytes"
	"io"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Load OpenAPI Specification v3.0.0 spec file.
func Load(filename string) (*Document, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	io.Copy(&buf, f)
	b := buf.Bytes()

	if err := f.Close(); err != nil {
		panic(err)
	}

	return parse(b)
}

func parse(b []byte) (*Document, error) {
	doc := Document{}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}
