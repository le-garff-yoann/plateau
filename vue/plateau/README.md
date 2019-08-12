# Plateau

[![Stability Status](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/orangemug/stability-badges)
[![Pipeline Status](https://gitlab.com/le-garff-yoann/plateau/badges/master/pipeline.svg)](https://gitlab.com/le-garff-yoann/plateau/pipelines)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

## Build

```bash
npm run build
```

## Other `run` subcommands

- Unit test

```bash
npm run test:unit
```

- Lint

```bash
npm run lint
```

- Server aka "development mode"

1. Run `plateau run` with the `-l :3000` parameter.
2. 
```bash
NODE_DEV_PROXY_API=http://localhost:3000 \
    npm run serve
```
