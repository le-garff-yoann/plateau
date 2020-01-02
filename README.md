# Plateau

<p align="center">
    <img src="vue/plateau/src/assets/logo.png" alt="Plateau" title="Plateau" />
</p>

[![Stability Status](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/orangemug/stability-badges)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Pipeline Status](https://gitlab.com/le-garff-yoann/plateau/badges/master/pipeline.svg)](https://gitlab.com/le-garff-yoann/plateau/pipelines)
![Go Coverage Report](https://gitlab.com/le-garff-yoann/plateau/badges/master/coverage.svg?job=go:unit%20tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/le-garff-yoann/plateau)](https://goreportcard.com/report/github.com/le-garff-yoann/plateau)
[![Docker Pull](https://img.shields.io/static/v1?label=docker%20pull&message=registry.gitlab.com/le-garff-yoann/plateau&color=informational)](https://registry.gitlab.com/le-garff-yoann/plateau)

> Build your own board game server. Batteries included!

## Basics - Build and run via Docker

The code in this repository will build the binary for the [Rock–paper–scissors](https://en.wikipedia.org/wiki/Rock%E2%80%93paper%E2%80%93scissors) game.

```bash
# Build the image using the inmemory store.
docker build . \
    -t my-plateau --build-arg GO_TAGS="run_rockpaperscissors run_inmemory"

# Run the server.
docker run -d -p 8080:80 my-plateau \
    run -l :80 --listen-static-dir /public --session-key my-STRONG-secret
```

## [Customizing](CUSTOMIZING.md)
