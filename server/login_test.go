package server

import (
	"net/http"
	"net/http/httptest"
	"plateau/store/inmemory"
	"strings"
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

func TestLoginHandlers(t *testing.T) {
	t.Parallel()

	srv, err := New("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	var (
		registerH = http.Handler(srv.router.Get("registerUser").GetHandler())
		loginH    = http.Handler(srv.router.Get("loginUser").GetHandler())
		logoutH   = http.Handler(srv.router.Get("logoutUser").GetHandler())
	)

	newRecorder := func(h http.Handler) *httptest.ResponseRecorder {
		req, err := http.NewRequest("POST", "", strings.NewReader(`{"username":"foo","password":"bar"}`))
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		return rr
	}

	require.Equal(t, 401, newRecorder(loginH).Code)
	require.Equal(t, 201, newRecorder(registerH).Code)
	require.Equal(t, 409, newRecorder(registerH).Code)
	require.Equal(t, 201, newRecorder(loginH).Code)
	require.Equal(t, 201, newRecorder(logoutH).Code)
}
