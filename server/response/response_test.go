package response

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/le-garff-yoann/plateau/server/response/body"
	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	var (
		reqBody = body.New()
		resBody *body.Body
	)

	w := httptest.NewRecorder()

	WriteJSON(w, http.StatusInternalServerError, reqBody)
	body, _ := ioutil.ReadAll(w.Result().Body)

	require.NoError(t, json.Unmarshal(body, &resBody))
	require.Equal(t, reqBody, resBody)
}
