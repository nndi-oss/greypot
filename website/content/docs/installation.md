+++
date = "2023-02-02-10T06:43:48+02:00"
title = "Installation"
draft = false
weight = 2
description = "Installing Greypot and Greypot Studio"
toc = true
bref = "Installation"
+++

### Docker

> NOTE: We are working on uploading images to Dockerhub - stay tuned.

### Build Docker Image from source code

```sh
$ git clone https://github.com/nndi-oss/greypot

$ cd greypot

$ docker build -t greypot-server .

$ docker run -p "7665:7665" -v "$(pwd)/examples/fiber_example/templates:/templates" greypot-server
```

### Prebuilt Binaries

> NOTE: We are working on this

### Build from Source

Clone the repo, and run the following

```sh
$ git clone https://github.com/nndi-oss/greypot

$ cd greypot

$ cd ui && npm install && npm run build

$ go build -o ./build/greypot-server cmd/greypot-server/*.go
```

### Playwright Requirements for Building from Source

Currently, we are focusing on making the playwright based renderer work really good! The base project used Chrome Developer Protocol to connect with a Chromium instance. We [decided](https://github.com/nndi-oss/greypot/issues/1) to remove support for that.

In order to use the [Playwright](https://github.com/playwright-community/playwright-go) rendering functionality, you will need to have the [playwright dependencies](https://playwright.dev/docs/cli#install-system-dependencies) installed.

Read [here](https://playwright.dev/docs/cli#install-system-dependencies) for more info. But in short, you can use the following command to do so:

```
$ npx playwright install-deps chromium
```


