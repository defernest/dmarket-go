package items_test

import (
	"encoding/json"
	"github.com/defernest/dmarket-go/dmarket"
	"github.com/defernest/dmarket-go/mocks/items"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func getAllItems(t *testing.T, tsURL string, query url.Values) []dmarket.Object {
	var objects []dmarket.Object
	for {
		req, err := http.NewRequest(http.MethodGet, tsURL+"/exchange/v1/market/items?"+query.Encode(), http.NoBody)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		var ims dmarket.GetItemsResponse
		err = json.Unmarshal(body, &ims)
		require.NoError(t, err)
		query.Set("cursor", ims.Cursor)
		if len(ims.Objects) == 0 {
			break
		}
		objects = append(objects, ims.Objects...)
	}
	return objects
}

func TestMustReturnSuccess(t *testing.T) {
	cases := []struct {
		name             string
		query            url.Values
		wantHTTPCode     int
		wantObjectsCount int
	}{
		{
			name: "success: common",
			query: map[string][]string{
				"gameId":   {"9a92"},
				"currency": {"USD"},
				"limit":    {"100"},
			},
			wantHTTPCode:     http.StatusOK,
			wantObjectsCount: 500,
		},
		{
			name: "success: objects count < limit",
			query: map[string][]string{
				"gameId":   {"9a92"},
				"currency": {"USD"},
				"limit":    {"100"},
			},
			wantHTTPCode:     http.StatusOK,
			wantObjectsCount: 1,
		},
		{
			name: "success: objects count == 0",
			query: map[string][]string{
				"gameId":   {"9a92"},
				"currency": {"USD"},
				"limit":    {"100"},
			},
			wantHTTPCode:     http.StatusOK,
			wantObjectsCount: 0,
		},
		{
			name: "success: correct limit",
			query: map[string][]string{
				"gameId":   {"9a92"},
				"currency": {"USD"},
				"limit":    {"100"},
			},
			wantHTTPCode:     http.StatusOK,
			wantObjectsCount: 110,
		},
		{
			name: "error: no gameId",
			query: map[string][]string{
				"currency": {"USD"},
				"limit":    {"100"},
			},
			wantHTTPCode:     http.StatusBadRequest,
			wantObjectsCount: 0,
		},
		{
			name: "error: no currency",
			query: map[string][]string{
				"gameId": {"9a92"},
				"limit":  {"100"},
			},
			wantHTTPCode:     http.StatusBadRequest,
			wantObjectsCount: 0,
		},
		{
			name: "error: currency != USD",
			query: map[string][]string{
				"gameId":   {"9a92"},
				"currency": {"RUB"},
				"limit":    {"100"},
			},
			wantHTTPCode:     http.StatusBadRequest,
			wantObjectsCount: 0,
		},
		{
			name: "error: limit less 0",
			query: map[string][]string{
				"gameId":   {"9a92"},
				"currency": {"USD"},
				"limit":    {"-1"},
			},
			wantHTTPCode:     http.StatusBadRequest,
			wantObjectsCount: 0,
		},
		{
			name: "success: limit == 0",
			query: map[string][]string{
				"gameId":   {"9a92"},
				"currency": {"USD"},
				"limit":    {"0"},
			},
			wantHTTPCode:     http.StatusBadRequest,
			wantObjectsCount: 0,
		},
		{
			name: "error: limit more 100",
			query: map[string][]string{
				"gameId":   {"9a92"},
				"currency": {"USD"},
				"limit":    {"101"},
			},
			wantHTTPCode:     http.StatusBadRequest,
			wantObjectsCount: 0,
		},
		{
			name: "error: priceFrom less 0",
			query: map[string][]string{
				"gameId":    {"9a92"},
				"currency":  {"USD"},
				"limit":     {"100"},
				"priceFrom": {"-1"},
			},
			wantHTTPCode:     http.StatusBadRequest,
			wantObjectsCount: 0,
		},
		{
			name: "error: priceFrom > priceTo",
			query: map[string][]string{
				"gameId":    {"9a92"},
				"currency":  {"USD"},
				"limit":     {"100"},
				"priceFrom": {"10"},
				"priceTo":   {"1"},
			},
			wantHTTPCode:     http.StatusBadRequest,
			wantObjectsCount: 0,
		},
		{
			name: "success: priceFrom == priceTo",
			query: map[string][]string{
				"gameId":    {"9a92"},
				"currency":  {"USD"},
				"limit":     {"100"},
				"priceFrom": {"10"},
				"priceTo":   {"10"},
			},
			wantHTTPCode:     http.StatusOK,
			wantObjectsCount: 10,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Handle(items.MustReturnSuccess(tc.wantObjectsCount).Endpoint())
			ts := httptest.NewServer(router)
			if tc.wantHTTPCode == http.StatusOK {
				objects := getAllItems(t, ts.URL, tc.query)
				require.Len(t, objects, tc.wantObjectsCount)
			} else {
				resp, err := http.Get(ts.URL + "/exchange/v1/market/items?" + tc.query.Encode())
				require.NoError(t, err)
				require.Equal(t, tc.wantHTTPCode, resp.StatusCode)
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				require.Equal(t, "400: Bad Request", string(body))
			}
			ts.Close()
		})
	}
	t.Run("query to exchange params", func(t *testing.T) {
		router := gin.New()
		router.Handle(items.MustReturnSuccess(100).Endpoint())
		ts := httptest.NewServer(router)

		priceFrom, priceTo := 1000, 2000
		query := url.Values{
			"gameId":    {"9a92"},
			"currency":  {"USD"},
			"limit":     {"100"},
			"title":     {"SomeTitle"},
			"priceFrom": {strconv.Itoa(priceFrom)},
			"priceTo":   {strconv.Itoa(priceTo)},
		}
		for i := 0; i < 1000; i++ {
			objects := getAllItems(t, ts.URL, query)
			for _, object := range objects {
				require.Equal(t, query.Get("gameId"), object.GameID)
				require.Equal(t, query.Get("gameId"), object.Extra.GameID)
				require.Equal(t, query.Get("title"), object.Title)
				require.Equal(t, query.Get("title"), object.Extra.Name)
				itemPrice, err := strconv.Atoi(object.Price.Usd)
				require.NoError(t, err)
				require.GreaterOrEqual(t, itemPrice, priceFrom)
				require.LessOrEqual(t, itemPrice, priceTo)
			}
		}
	})
}
