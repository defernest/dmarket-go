package tests_test

import (
	"context"
	"github.com/defernest/dmarket-go/dmarket"
	"github.com/defernest/dmarket-go/mocks"
	"github.com/defernest/dmarket-go/mocks/common"
	"github.com/defernest/dmarket-go/mocks/items"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetAllItems(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		itemscount := 1000
		ts := mocks.NewDmarketServer(items.MustReturnSuccess(itemscount))
		wg := sync.WaitGroup{}
		results := dmarket.NewExchange(ts.Client).Items.GetAllItemsFromDmarket(context.Background())
		var objects []dmarket.Object
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				r := <-results
				require.NoError(t, r.Error)
				if len(r.Objects) == 0 {
					return
				}
				objects = append(objects, r.Objects...)
			}
		}()
		wg.Wait()
		require.Len(t, objects, itemscount)
	})
	t.Run("cancel with context", func(t *testing.T) {
		ts := mocks.NewDmarketServer(items.MustReturnSuccess(5000))
		wg := sync.WaitGroup{}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		results := dmarket.NewExchange(ts.Client).Items.GetAllItemsFromDmarket(ctx)
		deadline, _ := ctx.Deadline()
		var objects []dmarket.Object
		require.Eventually(t, func() bool {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					r, open := <-results
					if !open || r.Error != nil {
						return
					}
					objects = append(objects, r.Objects...)
				}
			}()
			wg.Wait()
			return true
		}, time.Until(deadline.Add(500*time.Millisecond)), 10*time.Millisecond)
		require.NotZero(t, len(objects))
		cancel()
	})
}
func TestItems_GetItems(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		wantItems := 100
		ts := mocks.NewDmarketServer(items.MustReturnSuccess(wantItems))
		e := dmarket.NewExchange(ts.Client)
		response, err := e.Items.GetItems("/exchange/v1/market/items?")
		require.NoError(t, err)
		require.Len(t, response.Objects, wantItems)
	})
	t.Run("error: unmarshal error", func(t *testing.T) {
		ts := mocks.NewDmarketServer(common.MustReturnBadBody(http.MethodGet, "/exchange/v1/market/items"))
		e := dmarket.NewExchange(ts.Client)
		_, err := e.Items.GetItems("/exchange/v1/market/items?")
		require.ErrorIs(t, err, dmarket.ErrUnmarshalAPIResponse)
	})
	errTests := []struct {
		name    string
		errCode int
	}{
		{name: "StatusBadRequest", errCode: http.StatusBadRequest},
		{name: "StatusNotFound", errCode: http.StatusNotFound},
		{name: "StatusInternalServerError", errCode: http.StatusInternalServerError},
	}
	for _, tt := range errTests {
		t.Run(tt.name, func(t *testing.T) {
			ts := mocks.NewDmarketServer(common.MustReturnHTTPError(http.MethodGet, "/exchange/v1/market/items", tt.errCode))
			e := dmarket.NewExchange(ts.Client)
			_, err := e.Items.GetItems("/exchange/v1/market/items?")
			require.ErrorAs(t, err, &dmarket.ErrorRepresentation{})
		})
	}
}

func TestItems_GenerateItems(t *testing.T) {
	cases := []struct {
		name   string
		params items.Params
		count  int
	}{
		{
			name: "success",
			params: items.Params{
				GameId:    "9b92",
				Title:     "title",
				Currency:  "USD",
				PriceFrom: 0,
				PriceTo:   10,
			},
			count: 10,
		},
		{
			name: "success with 0 prices",
			params: items.Params{
				GameId:    "9b92",
				Title:     "title",
				Currency:  "USD",
				PriceFrom: 0,
				PriceTo:   0,
			},
			count: 10,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := tc.params.GenerateItems(tc.count)
			require.Len(t, is, tc.count)
			for _, object := range is {
				require.Equal(t, object.GameID, tc.params.GameId)
				require.Equal(t, object.Title, tc.params.Title)
				price, err := strconv.Atoi(object.Price.Usd)
				require.NoError(t, err)
				if tc.params.PriceFrom == 0 && tc.params.PriceTo == 0 {
					require.Positive(t, price)
				} else {
					require.GreaterOrEqual(t, price, tc.params.PriceFrom)
					require.LessOrEqual(t, price, tc.params.PriceTo)
				}
			}
		})
	}
}
