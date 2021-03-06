new_plateau() {
    if [[ -z $1  ]]
    then
        echo "Usage: $0 <game name>"

        return 1
    fi

    if [[ -d ~/$1 ]]
    then
        echo "~/$1/ already exists"

        return 1
    fi

    echo "Creates project $1..."

    echo "Copy $PWD to ~/$1/"
    cp -r $PWD ~/$1

    pushd ~/$1 || return 1

    cat > cmd/run_$1.go <<EOF
// +build run_$1

package cmd

import (
	"plateau/game/$1"
	"plateau/server"
)

func newGame() server.Game {
	return &$1.Game{}
}
EOF

    echo "Creates game/$1/"
    mkdir game/$1

    echo "Writes game/$1/$1.go"
    cat > game/$1/$1.go <<EOF
package $1

import (
	"fmt"
	"plateau/protocol"
)

// Game ...
type Game struct {
	name, description      string
	minPlayers, maxPlayers uint
}

// Init ...
func (s *Game) Init() error {
	s.name = "$1"

	s.description = ""

	s.minPlayers = 2
	s.maxPlayers = 2

	return nil
}

// Name ...
func (s *Game) Name() string {
	return s.name
}

// Description ...
func (s *Game) Description() string {
	return s.description
}

// IsMatchValid ...
func (s *Game) IsMatchValid(g *protocol.Match) error {
	return nil
}

// MinPlayers ...
func (s *Game) MinPlayers() uint {
	return s.minPlayers
}

// MaxPlayers ...
func (s *Game) MaxPlayers() uint {
	return s.maxPlayers
}
EOF

    echo "Writes game/$1/$1_test.go"
    cat > game/$1/$1_test.go <<EOF
package $1

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGame(t *testing.T) {
	t.Parallel()

	g := Game{}
	g.Init()

	require.Equal(t, "$1", g.Name())

	require.Equal(t, "", g.Description())

	require.Equal(t, uint(2), g.MinPlayers())
	require.Equal(t, uint(2), g.MaxPlayers())
}
EOF

    echo "Writes game/$1/context.go"
    cat > game/$1/context.go <<EOF
package $1

import (
	"plateau/protocol"
	"plateau/server"
	"plateau/store"
)

// Context ...
func (s *Game) Context(trn store.Transaction, reqContainer *protocol.RequestContainer) *server.Context {
	return server.NewContext()
}
EOF

    echo "Writes game/$1/context_test.go"
    cat > game/$1/context_test.go <<EOF
package $1

import (
	"plateau/protocol"
	"plateau/server"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGameRuntime(t *testing.T) {
	t.Parallel()

	testMatchRuntime := &server.TestMatchRuntime{
		T:           t,
		Game:        &Game{},
		Match:       protocol.Match{NumberOfPlayersRequired: 2},
		PlayersName: []string{"foo", "bar"},
	}

	server.SetupTestMatchRuntime(t, testMatchRuntime)
	defer func() {
	    require.NoError(t, testMatchRuntime.Stop())
	}()

	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToJoin, protocol.ResOK)
	testMatchRuntime.TestRequest("bar", protocol.ReqPlayerWantToJoin, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToStartTheMatch, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerAccepts, protocol.ResOK)
	testMatchRuntime.TestRequest("bar", protocol.ReqPlayerAccepts, protocol.ResOK)
}
EOF

    echo "Cleans project $1..."
    rm -Rf \
        .git/ \
        .gitlab-ci.yml \
        {.,vue/plateau}/*.md
        vendor/ \
        cmd/run_rockpaperscissors.go \
        game/rockpaperscissors

    cat > game/$1/context_test.go <<EOF
EOF

    echo 'Writes of README.md'
    cat > README.md <<EOF
# Plateau - $1

## Build

```bash
go build -tags="run_$1 run_inmemory" -o dist/plateau
```
EOF
    echo "Writes of $1/vue/plateau/README.md"
    cat > vue/plateau/README.md <<EOF
# Plateau - $1

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
EOF

    popd

    echo "Done"
}
