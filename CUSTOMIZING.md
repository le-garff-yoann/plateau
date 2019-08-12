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

### [Vue.js frontend](vue/plateau)

This frontend is "game-agnostic" but you can customize the [`Match`](vue/plateau/src/components/Core/Match.vue) component to better fit your game.

**N.B**: *Plateau* exposes a [REST API](server/server.go) so you can easily build another frontend from scratch.
