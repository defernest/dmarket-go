package tests

import (
	"encoding/hex"
	"github.com/defernest/dmarket-go/mocks"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/defernest/dmarket-go/dmarket"
	"github.com/defernest/dmarket-go/mocks/common"

	"github.com/stretchr/testify/require"
)

func TestDefaultClient_Do(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ts := mocks.NewDmarketServer(common.MustReturnStatusOK(http.MethodGet, "/"))
		apiClient, err := dmarket.NewClient(ts.URL(), ts.PublicKey, ts.PrivareKey)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)
		resp, err := apiClient.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, "200 OK", resp.Status)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, int64(2), resp.ContentLength)
		require.ElementsMatch(t, []byte("{}"), resp.Body.Bytes())
	})
	t.Run("error: request sign error", func(t *testing.T) {
		apiClient, err := dmarket.NewClient("client://localhost", "f7235e1a233478f20b60b6240c49afb4d5a9970eb2228cb44b4183047eb89�", "255df258b9252c04d29ae19a88ef8be2dc8d3654a90037b8881937b81vfcf87vf7235e1a236478f20b60b6240c49afb4d5a9970eb2228cb44b4123044eb89ec3")
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)
		_, err = apiClient.DefaultClient.Do(req)
		var hexerr hex.InvalidByteError
		require.ErrorAs(t, err, &hexerr)
	})
	t.Run("error: Do request error", func(t *testing.T) {
		ts := mocks.NewDmarketServer(common.MustReturnStatusOK(http.MethodGet, "/"))
		apiClient, err := dmarket.NewClient("client://localhost", ts.PublicKey, ts.PrivareKey)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)
		_, err = apiClient.DefaultClient.Do(req)
		var doerr *url.Error
		require.ErrorAs(t, err, &doerr)
	})
	t.Run("dont panic", func(t *testing.T) {
		ts := mocks.NewDmarketServer(common.MustReturnStatusOK(http.MethodGet, "/"))
		apiClient, err := dmarket.NewClient(ts.URL(), ts.PublicKey, ts.PrivareKey)
		require.NoError(t, err)
		_, err = apiClient.DefaultClient.Do(nil)
		require.Error(t, err)
	})
}

func Test_DefaultClient_Delete(t *testing.T) {
	ts := mocks.NewDmarketServer(common.MustReturnStatusOK(http.MethodDelete, "/"))
	apiClient, err := dmarket.NewClient(ts.URL(), ts.PublicKey, ts.PrivareKey)
	require.NoError(t, err)
	payload := strings.NewReader("{}")
	resp, err := apiClient.DefaultClient.Delete("/", payload)
	require.NoError(t, err)
	require.Equal(t, "{}", resp.Body.String())
}

func Test_DefaultClient_Get(t *testing.T) {
	ts := mocks.NewDmarketServer(common.MustReturnStatusOK(http.MethodGet, "/"))
	apiClient, err := dmarket.NewClient(ts.URL(), ts.PublicKey, ts.PrivareKey)
	require.NoError(t, err)
	resp, err := apiClient.DefaultClient.Get("/")
	require.NoError(t, err)
	require.Equal(t, "{}", resp.Body.String())
}

func Test_DefaultClient_Patch(t *testing.T) {
	ts := mocks.NewDmarketServer(common.MustReturnStatusOK(http.MethodPatch, "/"))
	apiClient, err := dmarket.NewClient(ts.URL(), ts.PublicKey, ts.PrivareKey)
	require.NoError(t, err)
	payload := strings.NewReader("{}")
	resp, err := apiClient.DefaultClient.Patch("/", payload)
	require.NoError(t, err)
	require.Equal(t, "{}", resp.Body.String())
}

func Test_DefaultClient_Post(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ts := mocks.NewDmarketServer(common.MustReturnStatusOK(http.MethodPost, "/"))
		apiClient, err := dmarket.NewClient(ts.URL(), ts.PublicKey, ts.PrivareKey)
		require.NoError(t, err)
		payload := strings.NewReader("{}")
		resp, err := apiClient.DefaultClient.Post("/", payload)
		require.NoError(t, err)
		require.Equal(t, "{}", resp.Body.String())
	})
}
