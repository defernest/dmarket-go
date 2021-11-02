package mocks_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/defernest/dmarket-go/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type testEndpoint struct{}

func (t testEndpoint) Endpoint() (httpMethod string, relativePath string, handler gin.HandlerFunc) {
	return "GET", "/", func(context *gin.Context) {
		context.String(http.StatusOK, "{}")
	}
}

func TestNewDmarketServer(t *testing.T) {
	ts := mocks.NewDmarketServer(testEndpoint{})
	resp, err := ts.Client.Get("/")
	defer ts.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDmarketServer_Close(t *testing.T) {
	ts := mocks.NewDmarketServer(testEndpoint{})
	ts.Close()
	_, err := ts.Client.Get("/")
	var urlerr *url.Error
	require.ErrorAs(t, err, &urlerr)
}
