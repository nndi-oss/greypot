//go:build ignore
// +build ignore

package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nndi-oss/greypot"
	"github.com/nndi-oss/greypot/examples"
	greypotGin "github.com/nndi-oss/greypot/http/gin"
)

func main() {
	app := gin.New()

	module := greypot.NewModule(
		greypot.WithRenderTimeout(10*time.Second),
		greypot.WithViewport(2048, 1920),
		greypot.WithDjangoTemplateEngine(),
		greypot.WithTemplatesFromFilesystem("./templates/"),
		greypot.WithPlaywrightRenderer(),
	)

	greypotGin.Use(app, module)

	embedModule := greypot.NewModule(
		greypot.WithRenderTimeout(10*time.Second),
		greypot.WithViewport(2048, 1920),
		greypot.WithTemplatesFromFS(examples.ExampleTemplatesFS),
		greypot.WithGolangTemplateEngine(),
		greypot.WithPlaywrightRenderer(),
	)

	greypotGin.Use(app.Group("/embedded/"), embedModule)

	app.Run(":3000")
}
