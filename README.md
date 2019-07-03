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

**Please read [these instructions](CUSTOMIZING.md) is you want to customize *plateau* for another game.**

```bash
# Build to use the process memory as the store.
go build -tags="run_rockpaperscissors run_inmemory" -o dist/plateau 
 
# plateau help # Print the global help.
# plateau help run # Print the help for the run subcommand.

# Start the server.
dist/plateau run -l :3000 --session-key my-STRONG-secret
```

**N.B.** Parameters to the `run` subcommand may vary function of the flags declared by `store.RunCommandSetter(*cobra.Command)` (and thus by the implementation of `store.Store`).

## A quick look at the API

### Got yourself a session

```bash
BASE=http://localhost:3000
COOKIE_NAME=plateau
COOKIE_FILE=me.cookie
USERINFO='{"username":"me","password":"1234"}'

# Register yourself.
curl $BASE/user/register -d $USERINFO

# Log in.
curl $BASE/user/login --cookie-jar $COOKIE_FILE -d $USERINFO
```

### Play!

```bash
# Create and returns a match.
match_id=$(
curl -b $COOKIE_FILE -X POST $BASE/api/matchs \
    -d '{"number_of_players_required":2}' \
    | jq -r .id
)

# Listen for changes.
curl -b $COOKIE_FILE \
    $BASE/api/matchs/$match_id/notifications &

# Read or modify the state of the match with an in-game request.
curl -b $COOKIE_FILE -X PATCH $BASE/api/matchs/$match_id \
    -d '{"request":"?"}'
```

## Test a 2-player game in a shell (t2pg)

```bash
. helpers.bash

t2pg_setupmatch P1 P2

t2pg_send P1 ? # t2pg_send P1 PLAY_ROCK
t2pg_send P2 ? # t2pg_send P1 PLAY_PAPER
```

## Frontend

Take a look [here](vue/plateau/).
