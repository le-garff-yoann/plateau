# Plateau

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
