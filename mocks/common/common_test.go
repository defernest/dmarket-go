package common_test

import (
	"bytes"
	"fmt"
	"github.com/defernest/dmarket-go/mocks/common"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestNewCommonEndpoint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var b = common.NewEndpointBehavior(http.MethodGet, "/", func(context *gin.Context) {
			context.Status(http.StatusOK)
		})
		router := gin.New()
		router.Handle(b.Endpoint())
		require.HTTPSuccess(t, router.ServeHTTP, http.MethodGet, "/", nil)
	})
}

func TestMustReturnStatusOK(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		body   *bytes.Buffer
	}{
		{name: "GET success", method: http.MethodGet, path: "/", body: bytes.NewBufferString("{}")},
		{name: "POST success", method: http.MethodPost, path: "/post", body: bytes.NewBufferString("{post:success}")},
		{name: "DELETE success", method: http.MethodDelete, path: "/item-delete", body: bytes.NewBufferString("{delete:success}")},
		{name: "PATCH success", method: http.MethodPatch, path: "/some-path", body: bytes.NewBufferString("{patch:success}")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Handle(common.MustReturnStatusOK(tt.method, tt.path).Endpoint())
			ts := httptest.NewServer(router)
			body := bytes.NewBuffer(tt.body.Bytes())
			req, err := http.NewRequest(tt.method, ts.URL+tt.path, body)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Equal(t, tt.body.Len(), int(resp.ContentLength))
			require.Equal(t, tt.body.Bytes(), respBody)
		})
	}
}

func TestMustReturnBadBody(t *testing.T) {
	router := gin.New()
	router.Handle(common.MustReturnBadBody(http.MethodGet, "/").Endpoint())
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "{", w.Body.String())
}

func TestMustReturnHTTPError(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			name        string
			method      string
			path        string
			wantHTTPErr int
		}{
			{name: "StatusNotFound success", method: http.MethodGet, path: "/", wantHTTPErr: http.StatusNotFound},
			{name: "StatusBadRequest success", method: http.MethodPost, path: "/post", wantHTTPErr: http.StatusBadRequest},
			{name: "StatusUnauthorized success", method: http.MethodDelete, path: "/item-delete", wantHTTPErr: http.StatusUnauthorized},
			{name: "StatusInternalServerError success", method: http.MethodPatch, path: "/some-path", wantHTTPErr: http.StatusInternalServerError},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.NotPanics(t, func() {
					common.MustReturnHTTPError(tt.method, tt.path, tt.wantHTTPErr)
				})
				router := gin.New()
				router.Handle(common.MustReturnHTTPError(tt.method, tt.path, tt.wantHTTPErr).Endpoint())
				w := httptest.NewRecorder()
				req, err := http.NewRequest(tt.method, tt.path, nil)
				require.NoError(t, err)
				router.ServeHTTP(w, req)
				require.Equal(t, tt.wantHTTPErr, w.Code)
				require.Contains(t, w.Body.String(), fmt.Sprintf("%d: %s", tt.wantHTTPErr, http.StatusText(tt.wantHTTPErr)))
			})
		}
	})
	t.Run("panics", func(t *testing.T) {
		tests := []struct {
			name    string
			errcode int
		}{
			{name: "error < 400", errcode: 200},
			{name: "error > 511", errcode: 555},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.Panics(t, func() {
					common.MustReturnHTTPError(http.MethodGet, "/", tt.errcode)
				})
			})
		}
	})
}

func TestMustDoLongResponse(t *testing.T) {
	router := gin.New()
	router.Handle(common.MustDoLongResponse(http.MethodGet, "/", 5).Endpoint())
	ts := httptest.NewServer(router)
	req, err := http.NewRequest(http.MethodGet, ts.URL, http.NoBody)
	require.NoError(t, err)
	client := &http.Client{Timeout: 3 * time.Second}
	_, err = client.Do(req)
	var urlerr *url.Error
	require.ErrorAs(t, err, &urlerr)
	require.True(t, err.(*url.Error).Timeout())
}

func TestMustReturnDoBadRepresentationError(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		code := http.StatusBadRequest
		router := gin.New()
		router.Handle(common.MustReturnDoBadRepresentationError(http.MethodGet, "/", code).Endpoint())
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)
		router.ServeHTTP(w, req)
		require.Equal(t, code, w.Code)
		require.Equal(t, "{", w.Body.String())
	})
	t.Run("panic", func(t *testing.T) {
		require.Panics(t, func() {
			common.MustReturnDoBadRepresentationError(http.MethodGet, "/", 200)
		})
	})
}
