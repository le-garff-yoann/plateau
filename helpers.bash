new_plateau() (
	set -e

    if [[ -z $1 ]]
    then
        echo "Usage: $0 <game name>"

        return 1
    fi

    if [[ -d ~/$1 ]]
    then
		# shellcheck disable=SC2088
        echo "~/$1/ already exists"

        return 1
    fi

    echo "Creating project $1..."

    echo "Copying $PWD to ~/$1/"
    cp -r "$PWD" "$HOME/$1"

    cd "$HOME/$1"

    cat > "cmd/run_$1.go" <<EOF
// +build run_$1

package cmd

import (
	"github.com/le-garff-yoann/plateau/game/$1"
	"github.com/le-garff-yoann/plateau/server"
)

func newGame() server.Game {
	return &$1.Game{}
}
EOF

    echo "Creating game/$1/"
    mkdir "game/$1"

    echo "Generating game/$1/$1.go"
    cat > "game/$1/$1.go" <<EOF
package $1

import (
	"fmt"

	"github.com/le-garff-yoann/plateau/protocol"
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

    echo "Generating game/$1/$1_test.go"
    cat > "game/$1/$1_test.go" <<EOF
package $1

import (
	"testing"

	"github.com/le-garff-yoann/plateau/plateau/protocol"
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

    echo "Generating game/$1/context.go"
    cat > "game/$1/context.go" <<EOF
package $1

import (
	"github.com/le-garff-yoann/plateau/plateau/protocol"
	"github.com/le-garff-yoann/plateau/plateau/server"
	"github.com/le-garff-yoann/plateau/plateau/store"
)

// Context ...
func (s *Game) Context(trn store.Transaction, reqContainer *protocol.RequestContainer) *server.Context {
	return server.NewContext()
}
EOF

    echo "Generating game/$1/context_test.go"
    cat > "game/$1/context_test.go" <<EOF
package $1

import (
	"testing"

	"github.com/le-garff-yoann/plateau/protocol"
	"github.com/le-garff-yoann/plateau/server"
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

    echo "Cleaning up project $1..."
    rm -Rf .git/ .gitlab-ci.yml {.,vue/plateau}/*.md \
        vendor/ cmd/run_rockpaperscissors.go game/rockpaperscissors

    cat > "game/$1/context_test.go" <<EOF
EOF

    echo "Generating README.md"
    cat > README.md <<EOF
# Plateau - $1

## Build

\`\`\`bash
go build -tags="run_$1 run_inmemory" -o dist/plateau
\`\`\`
EOF
    echo "Generating $1/vue/plateau/README.md"
    cat > vue/plateau/README.md <<EOF
# Plateau - $1

## Build

\`\`\`bash
npm run build
\`\`\`

## Other \`run\` subcommands

- Unit test

\`\`\`bash
npm run test:unit
\`\`\`

- Lint

\`\`\`bash
npm run lint
\`\`\`

- Server aka "development mode"

1. Run \`plateau run\` with the \`-l :3000\` parameter.
2. 
\`\`\`bash
NODE_DEV_PROXY_API=http://localhost:3000 \
    npm run serve
\`\`\`
EOF

    echo Done.
)
