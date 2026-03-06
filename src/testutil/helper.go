package testutil

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// DecodeResp decodes a JSON response body into T.
func DecodeResp[T any](t *testing.T, body io.Reader) T {
	t.Helper()
	var r T
	require.NoError(t, json.NewDecoder(body).Decode(&r))
	return r
}
