package mocks

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"
)

var limiter = rate.NewLimiter(5, 5)

type testEndpoint struct{}

func (t testEndpoint) Endpoint() (httpMethod string, relativePath string, handler gin.HandlerFunc) {
	return "GET", "/", func(context *gin.Context) {
		context.String(http.StatusOK, "{}")
	}
}

func testDo(req *http.Request) (*http.Response, error) {
	err := limiter.Wait(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("api: request rate limiter error: %w", err)
	}
	return http.DefaultClient.Do(req)
}

func Test_dmarketAuth(t *testing.T) {
	ts := NewDmarketServer(testEndpoint{})
	defer ts.Close()
	cases := []struct {
		name                                              string
		wrongTimestamp, wrongPath, wrongPublic, wrongSign bool
	}{
		{
			name:           "wrong time",
			wrongTimestamp: true,
			wrongPath:      false,
			wrongPublic:    false,
			wrongSign:      false,
		},
		{
			name:           "wrong msg",
			wrongTimestamp: false,
			wrongPath:      true,
			wrongPublic:    false,
			wrongSign:      false,
		},
		{
			name:           "wrong pub",
			wrongTimestamp: false,
			wrongPath:      false,
			wrongPublic:    true,
			wrongSign:      false,
		},
		{
			name:           "wrong sign",
			wrongTimestamp: false,
			wrongPath:      false,
			wrongPublic:    false,
			wrongSign:      true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ts.wrongGet(tc.wrongTimestamp, tc.wrongPath, tc.wrongPublic, tc.wrongSign)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, "401: Unauthorized", string(body))
		})
	}
}

func Test_checkHeaders(t *testing.T) {
	ts := NewDmarketServer(testEndpoint{})
	defer ts.Close()
	cases := []struct {
		name           string
		headers        http.Header
		wantHTTPCode   int
		wantBodyString string
	}{
		{
			name: "success",
			headers: map[string][]string{
				"Accept":         {"application/json"},
				"X-Api-Key":      {"7bf5f047bf9c7630854b21e34d2fb3765558ba22c3a0123825709eac65b36ab4"},
				"Content-Type":   {"application/json"},
				"X-Request-Sign": {"dmar ed25519 d5d114a06d282855378a3d47b16bdb1293c9bb1a44c7271792ea2953d9772fc669294bb0fbd5e563e23dc981fc0bd933543c1502e7c48e3d95e6035135a1320b"},
				"X-Sign-Date":    {"1633697260"},
			},
			wantHTTPCode:   http.StatusOK,
			wantBodyString: "{}",
		},
		{
			name: "public key error: wrong key len",
			headers: map[string][]string{
				"X-Api-Key": {"7bf5f047bf"},
			},
			wantHTTPCode:   http.StatusBadRequest,
			wantBodyString: "400: Bad Request",
		},
		{
			name: "public key error: decode error",
			headers: map[string][]string{
				"X-Api-Key": {"e519f24bf3189604ef0db025451cbbaae0ed32820ed8d5f91ea0c6e74d1a5cccc"},
			},
			wantHTTPCode:   http.StatusBadRequest,
			wantBodyString: "400: Bad Request",
		},
		{
			name: "sign error: empty sign",
			headers: map[string][]string{
				"X-Api-Key":      {"7bf5f047bf9c7630854b21e34d2fb3765558ba22c3a0123825709eac65b36ab4"},
				"X-Request-Sign": {""},
			},
			wantHTTPCode:   http.StatusBadRequest,
			wantBodyString: "400: Bad Request",
		},
		{
			name: "sign error: sign header without dmar",
			headers: map[string][]string{
				"X-Api-Key":      {"7bf5f047bf9c7630854b21e34d2fb3765558ba22c3a0123825709eac65b36ab4"},
				"X-Request-Sign": {"d5d114a06d282855378a3d47b16bdb1293c9bb1a44c7271792ea2953d9772fc669294bb0fbd5e563e23dc981fc0bd933543c1502e7c48e3d95e6035135a1320b"},
			},
			wantHTTPCode:   http.StatusBadRequest,
			wantBodyString: "400: Bad Request",
		},
		{
			name: "sign error: decode error",
			headers: map[string][]string{
				"X-Api-Key":      {"7bf5f047bf9c7630854b21e34d2fb3765558ba22c3a0123825709eac65b36ab4"},
				"X-Request-Sign": {"d5d114a06d282855378a3d47b16bdb1293c9bb1a44c7271792ea2953d9772fc669294bb0fbd5e563e23dc981fc0bd933543c1502e7c48e3d95e6035135a1320"},
			},
			wantHTTPCode:   http.StatusBadRequest,
			wantBodyString: "400: Bad Request",
		},
		{
			name: "sign date error: empty date",
			headers: map[string][]string{
				"X-Api-Key":      {"7bf5f047bf9c7630854b21e34d2fb3765558ba22c3a0123825709eac65b36ab4"},
				"X-Request-Sign": {"dmar ed25519 d5d114a06d282855378a3d47b16bdb1293c9bb1a44c7271792ea2953d9772fc669294bb0fbd5e563e23dc981fc0bd933543c1502e7c48e3d95e6035135a1320b"},
				"X-Sign-Date":    {""},
			},
			wantHTTPCode:   http.StatusBadRequest,
			wantBodyString: "400: Bad Request",
		},
		{
			name: "sign date error: decode error",
			headers: map[string][]string{
				"X-Api-Key":      {"7bf5f047bf9c7630854b21e34d2fb3765558ba22c3a0123825709eac65b36ab4"},
				"X-Request-Sign": {"dmar ed25519 d5d114a06d282855378a3d47b16bdb1293c9bb1a44c7271792ea2953d9772fc669294bb0fbd5e563e23dc981fc0bd933543c1502e7c48e3d95e6035135a1320b"},
				"X-Sign-Date":    {"163369726ERR"},
			},
			wantHTTPCode:   http.StatusBadRequest,
			wantBodyString: "400: Bad Request",
		},
		{
			name: "accept or content-type empty",
			headers: map[string][]string{
				"X-Api-Key":      {"7bf5f047bf9c7630854b21e34d2fb3765558ba22c3a0123825709eac65b36ab4"},
				"X-Request-Sign": {"dmar ed25519 d5d114a06d282855378a3d47b16bdb1293c9bb1a44c7271792ea2953d9772fc669294bb0fbd5e563e23dc981fc0bd933543c1502e7c48e3d95e6035135a1320b"},
				"X-Sign-Date":    {"1633697260"},
			},
			wantHTTPCode:   http.StatusBadRequest,
			wantBodyString: "400: Bad Request",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, ts.URL(), http.NoBody)
			require.NoError(t, err)
			req.Header = tc.headers
			resp, err := testDo(req)
			require.NoError(t, err)
			require.Equal(t, tc.wantHTTPCode, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tc.wantBodyString, string(body))
		})
	}
}

func Test_noRoute(t *testing.T) {
	ts := NewDmarketServer(testEndpoint{})
	defer ts.Close()
	req, err := http.NewRequest(http.MethodGet, ts.URL()+"/err-path", http.NoBody)
	require.NoError(t, err)
	resp, err := ts.Client.Do(req)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "404: Not Found", string(body))
}
