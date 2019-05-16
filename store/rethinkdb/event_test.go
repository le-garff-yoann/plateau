package rethinkdb

import (
	"plateau/store"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventContainerConversion(t *testing.T) {
	t.Parallel()

	ec := EventContainer{}
	require.IsType(t, &store.EventContainer{}, ec.toStoreStruct())
	require.IsType(t, ec, *eventContainerFromStoreStruct(*ec.toStoreStruct()))
}
