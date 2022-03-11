package main

import (
	"net/http"

	"github.com/GuyARoss/orbit/examples/basic-react/orbitgen"
)

func main() {
	orb, err := orbitgen.New()
	if err != nil {
		panic(err)
	}

	orb.HandleFunc("/", func(c *orbitgen.Request) {
		c.RenderPage(orbitgen.ExampleTwoPage, nil)
	})

	http.ListenAndServe(":3030", orb.Serve())
}
