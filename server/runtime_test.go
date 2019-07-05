package server

import (
	"plateau/protocol"
	"testing"
)

func TestGameRuntime(t *testing.T) {
	t.Parallel()

	testMatchRuntime := &TestMatchRuntime{
		T:           t,
		Game:        &surrenderGame{},
		Match:       protocol.Match{NumberOfPlayersRequired: 2},
		PlayersName: []string{"foo", "bar"},
	}

	SetupTestMatchRuntime(t, testMatchRuntime)

	testMatchRuntime.TestRequest("foo", protocol.ReqListRequests, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToJoin, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToStartTheMatch, protocol.ResForbidden)
	testMatchRuntime.TestRequest("bar", protocol.ReqPlayerWantToJoin, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToLeave, protocol.ResOK)
	testMatchRuntime.TestRequest("bar", protocol.ReqPlayerWantToStartTheMatch, protocol.ResForbidden)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToJoin, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToStartTheMatch, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerAccepts, protocol.ResOK)
	testMatchRuntime.TestRequest("bar", protocol.ReqPlayerAccepts, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqListRequests, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerAccepts, protocol.ResNotImplemented)
}
