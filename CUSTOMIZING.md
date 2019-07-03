# Customizing

## Create a new project

```bash
. helpers.bash

set -e

new_plateau mygame # Create ~/mygame/
```

## It's now yours!

### Let's code

Your game logic should preferably stand in `~/mygame/game/mygame/`. It must implement the [`plateau/server.Game`](server/game.go) interface.

### Useful godoc

```bash
go doc plateau/server.Game
go doc plateau/server.TestMatchRuntime
go doc plateau/server.Context

go doc plateau/store.Transaction

go doc -all plateau/protocol
```
