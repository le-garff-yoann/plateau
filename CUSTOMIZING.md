# Customizing

## Create a new project

```bash
. helpers.bash

new_plateau mygame # Create ~/mygame/
```

## It's now yours to modify!

Your game logic should preferably stand in `~/mygame/game/mygame/`.
It must implement the [`plateau/server.Game`](server/game.go) interface.

### Useful godoc

- `plateau/server.Game`
- `plateau/server.TestMatchRuntime`
- `plateau/server.Context`
- `plateau/store.Transaction`
- `plateau/protocol`
