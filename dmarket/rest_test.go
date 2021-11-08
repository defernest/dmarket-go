package dmarket

import (
	"bytes"
	"errors"
	"net/http"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

func TestErrorRepresentation_Error(t *testing.T) {
	err := ErrorRepresentation{
		Response: Response{StatusCode: http.StatusUnauthorized},
	}
	require.Contains(t, err.Error(), "401")
	require.Contains(t, err.Error(), "Unauthorized")
}

func TestRuntimeError_Error(t *testing.T) {
	t.Run("runtimeError Error", func(t *testing.T) {
		want := "dmarket API runtime error::: err: runtime [code: 500], message: runtimeError error details: [0](typeURL: TypeURL, value: https://url.url/path?q=q) "
		err := RuntimeError{
			Err:     "runtime",
			Code:    500,
			Message: "runtimeError error",
			Details: []struct {
				TypeURL string "json:\"type_url\""
				Value   string "json:\"value\""
			}{{TypeURL: "TypeURL", Value: "https://url.url/path?q=q"}},
		}
		require.Equal(t, want, err.Error())
	})
}

func TestResponse_ReadFrom(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var r Response
		test := []byte{1, 2, 3, 4}
		i, err := r.ReadFrom(bytes.NewBuffer(test))
		require.NoError(t, err)
		require.Equal(t, int64(4), i)
		require.Equal(t, test, r.Body.Bytes())
	})
	t.Run("error", func(t *testing.T) {
		var r Response
		i, err := r.ReadFrom(iotest.ErrReader(errors.New("reader error")))
		require.Error(t, err)
		require.Zero(t, i)
	})
}
