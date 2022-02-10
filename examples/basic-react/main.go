package main

import (
	"time"

	"github.com/GuyARoss/orbit/examples/basic-react/orbitgen"
)

func main() {
	orbitgen.HandleFunc("/", func(c *orbitgen.RuntimeCtx) {
		now := time.Now()

		props := make(map[string]interface{})
		props["day"] = now.Day()
		props["month"] = now.Month()
		props["year"] = now.Year()

		c.RenderPage(orbitgen.ExamplePage, props)
	})

	orbitgen.Start(3003)
}
