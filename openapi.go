package openapi

import (
	"bytes"
	"io"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// LoadFile OpenAPI Specification v3.0 spec file.
func LoadFile(filename string) (*Document, error) {
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

	return Load(b)
}

// Load OpenAPI Specification v3.0 spec.
func Load(b []byte) (*Document, error) {
	doc := &Document{}
	if err := yaml.Unmarshal(b, doc); err != nil {
		return nil, err
	}
	for i := range doc.Security {
		doc.Security[i].setDocument(doc)
	}
	return doc, nil
}
