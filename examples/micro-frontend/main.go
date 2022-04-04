package main

import (
	"net/http"

	"github.com/GuyARoss/orbit/examples/micro-frontend/orbitgen"
)

func main() {
	orb, err := orbitgen.New()
	if err != nil {
		panic(err)
	}

	orb.HandleFunc("/", func(c *orbitgen.Request) {
		props := make(map[string]interface{})
		props["age"] = 23
		props["name"] = "Guy"

		c.RenderPages(props, orbitgen.NamePage, orbitgen.AgePage, orbitgen.StaticPage)
	})

	http.ListenAndServe(":3030", orb.Serve())
}
