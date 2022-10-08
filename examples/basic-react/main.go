package main

import (
	"net/http"
	"time"

	"github.com/GuyARoss/orbit/examples/basic-react/orbitgen"
)

func main() {
	orb, err := orbitgen.New()
	if err != nil {
		panic(err)
	}

	orb.HandleFunc("/", func(c *orbitgen.Request) {
		now := time.Now()

		props := make(map[string]interface{})
		props["day"] = now.Day()
		props["month"] = now.Month()
		props["year"] = now.Year()

		c.RenderPage(orbitgen.ExamplePage, props)
	})

	orb.HandleFunc("/second", func(c *orbitgen.Request) {
		now := time.Now()

		props := make(map[string]interface{})
		props["day"] = now.Day()
		props["month"] = now.Month()
		props["year"] = now.Year()

		c.RenderPage(orbitgen.ExampleTwoPage, props)
	})

	err = http.ListenAndServe(":3030", orb.Serve())
	if err != nil {
		panic(err)
	}
}
