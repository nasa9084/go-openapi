package main

import (
	"log"

	"github.com/nasa9084/go-openapi/internal/astutil"
	"github.com/nasa9084/go-openapi/internal/generators/getter"
	"github.com/nasa9084/go-openapi/internal/generators/resolve"
	"github.com/nasa9084/go-openapi/internal/generators/setroot"
	"github.com/nasa9084/go-openapi/internal/generators/unmarshalyaml"

	"golang.org/x/sync/errgroup"
)

type Generator interface {
	Generate() error
}

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	objects, err := astutil.ParseOpenAPIObjects("interfaces.go")
	if err != nil {
		return err
	}

	generators := []Generator{
		unmarshalyaml.NewGenerator(objects),
		resolve.NewGenerator(objects),
		getter.NewGenerator(objects),
		setroot.NewGenerator(objects),
	}

	var eg errgroup.Group
	for _, generator := range generators {
		eg.Go(generator.Generate)
	}

	return eg.Wait()
}
