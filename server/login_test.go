package server

import (
	"fmt"
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

func testRegisterAndLoginHandlers(t *testing.T, srv *Server, username, password string) (func(h http.Handler) *httptest.ResponseRecorder, *httptest.ResponseRecorder) {
	var (
		registerH = http.Handler(srv.router.Get("registerUser").GetHandler())
		loginH    = http.Handler(srv.router.Get("loginUser").GetHandler())
	)

	newRecorder := func(h http.Handler) *httptest.ResponseRecorder {
		req, err := http.NewRequest("POST", "", strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		return rr
	}

	require.Equal(t, http.StatusUnauthorized, newRecorder(loginH).Code)
	require.Equal(t, http.StatusCreated, newRecorder(registerH).Code)
	require.Equal(t, http.StatusConflict, newRecorder(registerH).Code)

	loginRecorder := newRecorder(loginH)
	require.Equal(t, http.StatusCreated, loginRecorder.Code)

	return newRecorder, loginRecorder
}

func TestRegisterLoginAndLogoutHandlers(t *testing.T) {
	t.Parallel()

	srv, err := Init("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	newRecorder, _ := testRegisterAndLoginHandlers(t, srv, "foo", "bar")

	logoutH := http.Handler(srv.router.Get("logoutUser").GetHandler())

	require.Equal(t, http.StatusCreated, newRecorder(logoutH).Code)
}
