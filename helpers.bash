new_plateau() {
    if [[ -z $1  ]]
    then
        echo "Usage: $0 <game name>"

        return 1
    fi

    echo "Copy $PWD to ~/$1"
    cp -r $PWD ~/$1

    pushd ~/$1

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

    echo "Create game/$1/"
    mkdir game/$1

    echo "Write game/$1/$1.go"
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
    echo "Write game/$1/$1_test.go"
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

    echo "Write game/$1/context.go"
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
    echo "Write game/$1/context_test.go"
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

	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToJoin, protocol.ResOK)
	testMatchRuntime.TestRequest("bar", protocol.ReqPlayerWantToJoin, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToStartTheGame, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerAccepts, protocol.ResOK)
	testMatchRuntime.TestRequest("bar", protocol.ReqPlayerAccepts, protocol.ResOK)
}
EOF

    echo "Cleaning the project"
    rm -Rf \
        vendor/ \
        .git/ \
        cmd/run_rockpaperscissors.go \
        game/rockpaperscissors

    popd
}

_t2pg_req() {
    [[ ! $(which curl jq) || -z $1 ]] && return 1 

    local \
        BASE=${T2PG_PLATEAU_BASEURL:-http://localhost:3000} \
        COOKIE_NAME=plateau \
        COOKIE_FILE=$1.cookie \
        USERINFO="{\"username\":\"$1\",\"password\":\"$1\"}"

    curl $BASE/user/register -d $USERINFO &>/dev/null
    curl $BASE/user/login --cookie-jar $COOKIE_FILE -d $USERINFO 2>/dev/null

    local match_id=$(curl -b $COOKIE_FILE $BASE/api/matchs 2>/dev/null | jq -r '.[0]')

    [[ $match_id == "null" ]] && \
    match_id=$(curl 2>/dev/null -b $COOKIE_FILE -X POST $BASE/api/matchs \
        -d '{"number_of_players_required":2}' | jq -r .id)

    curl $BASE/api/matchs/$match_id$2 -b $COOKIE_FILE ${@:3} 2>/dev/null | jq .

    [[ ${PIPESTATUS[0]} -eq 0 ]]
}

t2pg_cleanup() {
    [[ -z $1 ]] && return 1

    rm -f $1.cookie
}

t2pg_match() {
    [[ -z $1 ]] && return 1

    _t2pg_req $1 /
}

t2pg_deals() {
    [[ -z $1 ]] && return 1

    _t2pg_req $1 /deals
}

t2pg_send() {
    [[ -z $1 || -z $2 ]] && return 1

    _t2pg_req $1 / -X PATCH -d "{\"request\":\"$2\"}"
}

t2pg_setupmatch() {
    [[ -z $1 || -z $2 ]] && return 1 

    (
        t2pg_send $1 PLAYER_WANT_TO_JOIN && \
        t2pg_send $2 PLAYER_WANT_TO_JOIN && \
        t2pg_send $2 PLAYER_WANT_TO_START_THE_GAME && \
        t2pg_send $2 PLAYER_ACCEPTS && \
        t2pg_send $1 PLAYER_ACCEPTS
    ) 1>/dev/null && echo "Done."
}
