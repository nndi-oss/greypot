//go:build ignore
// +build ignore

package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot"
	"github.com/nndi-oss/greypot/examples"
	greypotFiber "github.com/nndi-oss/greypot/http/fiber"
)

func main() {
	app := fiber.New()

	module := greypot.NewModule(
		greypot.WithRenderTimeout(10*time.Second),
		greypot.WithViewport(2048, 1920),
		greypot.WithDjangoTemplateEngine(),
		greypot.WithTemplatesFromFilesystem("./templates/"),
		greypot.WithPlaywrightRenderer(),
	)

	greypotFiber.Use(app, module)

	embedModule := greypot.NewModule(
		greypot.WithRenderTimeout(10*time.Second),
		greypot.WithViewport(2048, 1920),
		greypot.WithTemplatesFromFS(examples.ExampleTemplatesFS),
		greypot.WithGolangTemplateEngine(),
		greypot.WithPlaywrightRenderer(),
	)

	greypotFiber.Use(app.Group("/embedded/"), embedModule)

	app.Listen(":3000")
}
