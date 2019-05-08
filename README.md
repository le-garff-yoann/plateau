# Plateau

<p align="center">
    <img src="docs/content/assets/img/plateau-logo.png" alt="Plateau" title="Plateau" />
</p>

[![Stability Status](https://img.shields.io/badge/stability-work_in_progress-red.svg)](https://github.com/orangemug/stability-badges)
[![Pipeline Status](https://gitlab.com/le-garff-yoann/plateau/badges/master/pipeline.svg)](https://gitlab.com/le-garff-yoann/plateau/pipelines)
![Go Coverage Report](https://gitlab.com/le-garff-yoann/plateau/badges/master/coverage.svg?job=go:test)
[![Go Report Card](https://goreportcard.com/badge/github.com/le-garff-yoann/plateau)](https://goreportcard.com/report/github.com/le-garff-yoann/plateau)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

> **The code in this repository will build the binary for the [Rock–paper–scissors](https://en.wikipedia.org/wiki/Rock%E2%80%93paper%E2%80%93scissors)** game.
Please read [these instructions](CUSTOMIZING.md) is you want to "customize and build" for another game.

## Run

```bash
# plateau help # Print the global help.
# plateau help run # Print the help for the run subcommand.

# Start the server.
plateau run \
    --listen :3000 \
    --listen-session-key my-STRONG-secret \
    --pg-conn-str postgres://pg:pg@localhost/pg
```

## API

### Got yourself a session

```bash
BASE=http://localhost:3000
COOKIE_NAME=plateau
PLAYER_NAME=me
PLAYER_PASSWORD=me1234

# Register yourself.
curl $BASE/user/register \
    -d "{\"username\":\"$PLAYER_NAME\",\"password\":\"$PLAYER_PASSWORD\"}"

# Log in.
curl $BASE/user/login --cookie-jar cookies.out \
    -d "{\"username\":\"$PLAYER_NAME\",\"password\":\"$PLAYER_PASSWORD\"}"
```

### Play!

```bash
# Returns players.
curl -b cookies.out $BASE/api/players

# Create and returns a game.
game=$(
curl -b cookies.out -X POST $BASE/api/games \
    -d '{"number_of_players_required":2}'
)

# Connect to the game throught WebSocket.
wscat \
    -H X-Interactive:true \
    -H Cookie:$COOKIE_NAME=$(grep $COOKIE_NAME cookies.out | awk '{print $7}') \
    -c $BASE/api/games/$(echo $game | jq .id)
```

## Frontend

Take a look [here](vue/plateau/).
