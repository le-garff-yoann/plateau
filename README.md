# Plateau

<p align="center">
    <img src="docs/content/assets/img/plateau-logo.png" alt="Plateau" title="Plateau" />
</p>

[![Stability Status](https://img.shields.io/badge/stability-work_in_progress-red.svg)](https://github.com/orangemug/stability-badges)
[![Pipeline Status](https://gitlab.com/le-garff-yoann/plateau/badges/master/pipeline.svg)](https://gitlab.com/le-garff-yoann/plateau/pipelines)
![Go Coverage Report](https://gitlab.com/le-garff-yoann/plateau/badges/master/coverage.svg?job=go:test)
[![Go Report Card](https://goreportcard.com/badge/github.com/le-garff-yoann/plateau)](https://goreportcard.com/report/github.com/le-garff-yoann/plateau)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

> Build your own board game server. Batteries included!

## Let's go

**The code in this repository will build the binary for the [Rock–paper–scissors](https://en.wikipedia.org/wiki/Rock%E2%80%93paper%E2%80%93scissors) game.**

```bash
go build -tags="game_rockpaperscissors run_rethinkdb" -o plateau # And build to use RethinkDB as the store.

# ./plateau help # Print the global help.
# ./plateau help run # Print the help for the run subcommand.

# Start the server.
./plateau run \
    --listen :3000 \
    --session-key my-STRONG-secret \
    --rethinkdb-address rethinkdb:28015 \
    --rethinkdb-database plateau \
    --rethinkdb-create-tables
```

**N.B.** Please read [these instructions](CUSTOMIZING.md) is you want to "customize and build" for another game.

**N.B.** Parameters to the `run` subcommand may vary function of the flags declared by `store.RunCommandSetter(*cobra.Command)` (and thus by the implementation of `store.Store`).

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

# Create and returns a match.
match=$(
curl -b cookies.out -X POST $BASE/api/matchs \
    -d '{"number_of_players_required":2}'
)

# Connect to the match throught WebSocket.
wscat \
    -H X-Interactive:true \
    -H Cookie:$COOKIE_NAME=$(grep $COOKIE_NAME cookies.out | awk '{print $7}') \
    -c $BASE/api/matchs/$(echo $match | jq -r .id)
```

## Frontend

Take a look [here](vue/plateau/).
