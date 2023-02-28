# greypot

An experimental library for generating PDF reports from HTML templates or "designs".

You can use it in three ways

* As a Go Library
* As an API which can be interacted with from any library
* Use the Greypot Studio for designing and testing the templates in browser

## Quick API Example

You can use the following command to generate a PDF using the Greypot Studio API running on Fly.io

```sh
$ curl --request POST --url https://greypot-studio.fly.dev/_studio/generate/pdf/test \
  --header 'Content-Type: application/json' \
  --data '{ "Name": "test.html", "Template": "<h1>Hello {{data.name}}</h1>", "Data": { "name": "John Smith" } }' | jq -r '.data' | base64 --decode > test.pdf
```

## Use as a Go Library

Say you want to produce reports or other such type of documents in your applications. 
`greypot` allows you to design your reports with HTML as template files that use  a Django-like [templating](https://docs.djangoproject.com/en/4.1/ref/templates/language/) engine. We also support the standard Go `html/template`.

These HTML reports can then be generated as HTML, PNG or PDF via endpoints that greypot adds to your application when you use the framework support (for Fiber or Gin).

Once you add the middleware to your application, it adds the following routes:

```
GET /reports/list

GET /reports/preview/:reportTemplateName

GET /reports/render/:reportTemplateName

POST /reports/export/html/:reportTemplateName

POST /reports/export/png/:reportTemplateName

POST /reports/export/pdf/:reportTemplateName

POST /reports/export/bulk/html/:reportTemplateName

POST /reports/export/bulk/png/:reportTemplateName

POST /reports/export/bulk/pdf/:reportTemplateName
```

You can then call these from within your applications to generate/export the reports e.g. from a frontend UI.

### Using Greypot with Gin

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


### Using Greypot with Fiber


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


### Build Docker Image and run the Container

```sh
$ docker build -t greypot-server .

$ docker run -p "7665:7665" -v "$(pwd)/examples/fiber_example/templates:/templates" greypot-server
```


## Playwright Rendering Engine

Currently, we are focusing on making the playwright based renderer work really good! The base project used Chrome Developer Protocol to connect with a Chromium instance. We [decided](https://github.com/nndi-oss/greypot/issues/1) to remove support for that.

In order to use the [Playwright](https://github.com/playwright-community/playwright-go) rendering functionality, you will need to have the [playwright dependencies](https://playwright.dev/docs/cli#install-system-dependencies) installed.

Read [here](https://playwright.dev/docs/cli#install-system-dependencies) for more info. But in short, you can use the following command to do so:

```
$ npx playwright install-deps chromium
```

## Prior art

- Greypot is based on [go-report-builder](https://github.com/AdikaStyle/go-report-builder)
