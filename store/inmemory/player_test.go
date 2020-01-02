package inmemory

import (
	"testing"
)

func TestPlayerList(t *testing.T) {
	testStr(t).TestPlayerList()
}

func TestPlayerCreate(t *testing.T) {
	testStr(t).TestPlayerCreate()
}

func TestPlayerRead(t *testing.T) {
	testStr(t).TestPlayerRead()
}

func TestPlayerIncreaseWins(t *testing.T) {
	testStr(t).TestPlayerIncreaseWins()
}

func TestPlayerIncreaseLoses(t *testing.T) {
	testStr(t).TestPlayerIncreaseLoses()
}

func TestPlayerIncreaseTies(t *testing.T) {
	testStr(t).TestPlayerIncreaseTies()
}
