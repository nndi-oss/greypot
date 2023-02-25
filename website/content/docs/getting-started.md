+++
date = "2023-02-02-10T06:43:48+02:00"
title = "Getting Started"
draft = false
weight = 1
description = "Getting Started with Greypot and Greypot Studio"
toc = true
bref = "GettingStarted"
+++

Greypot converts your HTML into PDF, HTML, and PNG. You can use it in three ways

* As a Go Library
* As an API which can be interacted with from any library
* Use the Greypot Studio for designing and testing the templates in browser


### Use Greypot Studio Cloud

Use the API on our [Greypot Cloud](https://greypot-studio.fly.dev) Service - which includes the Studio UI for designing and prototyping your report designs.

Read the [API Documentation here]({{< ref "openapi" >}})



### Use as a Go Library

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

See [the tutorial]({{< ref "gin-support" >}})

### Using Greypot with Fiber

See [the tutorial]({{< ref "gofiber-support" >}})

### Use Greypot Studio from Docker

See [the Docker guide]({{< ref "installation.md#docker" >}})

