+++
date = "2023-02-02-10T06:43:48+02:00"
title = "Greypot: Fiber Framework Support for Go developers"
draft = false
weight = 10
description = "Fiber Framework Support for Go developers"
toc = true
bref = "Installation"
+++


### Basic Usage


```go
package main

import (
	"embed"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot"
	greypotFiber "github.com/nndi-oss/greypot/http/fiber"
)

//go:embed "templates"
var templatesFS embed.FS

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
		greypot.WithTemplatesFromFS(templatesFS),
		greypot.WithGolangTemplateEngine(),
		greypot.WithPlaywrightRenderer(),
	)

	greypotFiber.Use(app.Group("/embedded/"), embedModule)

	app.Listen(":3000")
}
```