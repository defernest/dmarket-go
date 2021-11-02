package dmarket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItemsPriceRange(t *testing.T) {
	type args struct {
		priceFrom int
		priceTo   int
	}
	t.Run("price range setup success", func(t *testing.T) {
		tests := []struct {
			name string
			e    Items
			args args
			err  error
		}{
			{name: "OK", e: Items{}, args: args{priceFrom: 1, priceTo: 2}, err: nil},
			{name: "OK:priceFrom==priceTo", e: Items{}, args: args{priceFrom: 1, priceTo: 1}, err: nil},
			{name: "ERR:priceFrom<priceTo", e: Items{}, args: args{2, 1}, err: ErrIncorrectPriceRange},
			{name: "ERR:priceFrom<0", e: Items{}, args: args{-1, 1}, err: ErrIncorrectPriceRange},
			{name: "ERR:priceTo<0", e: Items{}, args: args{1, -1}, err: ErrIncorrectPriceRange},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.err != nil {
					require.Panics(t, func() {
						ItemsPriceRange(tt.args.priceFrom, tt.args.priceTo)(&tt.e)
					})
				} else {
					ItemsPriceRange(tt.args.priceFrom, tt.args.priceTo)(&tt.e)
					require.Equal(t, tt.args.priceTo, tt.e.priceTo)
					require.Equal(t, tt.args.priceFrom, tt.e.priceFrom)
				}
			})
		}
	})
}

func TestItemsLimitPerRequest(t *testing.T) {
	t.Run("item limit per request success", func(t *testing.T) {
		tests := []struct {
			name  string
			e     Items
			limit int
			err   error
		}{
			{name: "OK:limit==10", e: Items{}, limit: 10, err: nil},
			{name: "ERR:limit==0", e: Items{}, limit: 0, err: ErrLimitPerRequest},
			{name: "ERR:limit<=0", e: Items{}, limit: -1, err: ErrLimitPerRequest},
			{name: "ERR:limit>100", e: Items{}, limit: 101, err: ErrLimitPerRequest},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.err != nil {
					require.Panics(t, func() {
						ItemsLimitPerRequest(tt.limit)(&tt.e)
					})
				} else {
					ItemsLimitPerRequest(tt.limit)(&tt.e)
					require.Equal(t, tt.limit, tt.e.limit)
				}
			})
		}
	})
}

func TestItemsTitle(t *testing.T) {
	t.Run("success: title", func(t *testing.T) {
		title := "test"
		i := Items{}
		ItemsTitle(title)(&i)
		require.Equal(t, title, i.title)
	})
}
