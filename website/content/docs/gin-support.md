+++
date = "2023-02-02-10T06:43:48+02:00"
title = "Greypot: Gin Framework Support for Go developers"
draft = false
weight = 10
description = "Gin Framework Support for Go developers"
toc = true
bref = "Installation"
+++

### Basic Usage

```go
package main

import (
	"embed"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nndi-oss/greypot"
	greypotGin "github.com/nndi-oss/greypot/http/gin"
)

//go:embed "templates"
var templatesFS embed.FS

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
		greypot.WithTemplatesFromFS(templatesFS),
		greypot.WithGolangTemplateEngine(),
		greypot.WithPlaywrightRenderer(),
	)

	greypotGin.Use(app.Group("/embedded/"), embedModule)

	app.Run(":3000")
}
```
