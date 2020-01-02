package postgresql

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchConversion(t *testing.T) {
	t.Parallel()

	match := &Match{}

	require.IsType(t, &protocol.Match{}, match.toProtocolStruct())
	require.IsType(t, match, matchFromProtocolStruct(match.toProtocolStruct()))
}

func TestMatchCreateAndList(t *testing.T) {
	testStr(t).TestMatchCreateAndList()
}

func TestMatchRead(t *testing.T) {
	testStr(t).TestMatchRead()
}

func TestMatchEndedAt(t *testing.T) {
	testStr(t).TestMatchEndedAt()
}

func TestMatchCreateDeal(t *testing.T) {
	testStr(t).TestMatchCreateDeal()
}

func TestMatchUpdateCurrentDealHolder(t *testing.T) {
	testStr(t).TestMatchUpdateCurrentDealHolder()
}

func TestMatchPlayerJoins(t *testing.T) {
	testStr(t).TestMatchPlayerJoins()
}

func TestMatchPlayerLeaves(t *testing.T) {
	testStr(t).TestMatchPlayerLeaves()
}
