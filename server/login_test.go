package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidLoginCredentials(t *testing.T) {
	t.Parallel()

	cred := loginCredentials{Username: "foo", Password: "bar"}

	require.True(t, cred.IsValid())
}

func TestInvalidLoginCredentials(t *testing.T) {
	t.Parallel()

	cred := loginCredentials{Username: "", Password: ""}

	require.False(t, cred.IsValid())
}

func TestLoginCredentialsHashedPassword(t *testing.T) {
	t.Parallel()

	cred := loginCredentials{Username: "foo", Password: "bar"}

	hPassword, err := cred.PasswordHash()
	require.NoError(t, err)
	require.True(t, cred.VerifyHash(hPassword))
}
