package dmarket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	marketItems = "/exchange/v1/market/items?"
	userItems   = "/exchange/v1/user/items?"
)

var (
	// ErrUnexpectedAPIResponse
	ErrUnexpectedAPIResponse = errors.New("unexpected HTTP error from dmarket items endpoint")
	// ErrUnmarshalAPIResponse
	ErrUnmarshalAPIResponse = errors.New("can not unmarshal resp json")
	// ErrIncorrectPriceRange displays an error of an incorrectly defined price range for request
	ErrIncorrectPriceRange = errors.New("incorrect price range")
	// ErrLimitPerRequest indicates an error when setting the item limit for a request
	ErrLimitPerRequest = errors.New("item limit for the request cannot be less than or equal to zero")
)

//Items is a service structure for interacting with dmarket Items API endpoint
type Items struct {
	client                    Requester
	title, cursor             string
	priceFrom, priceTo, limit int
}

// Options is functional option for Items endpoint
type Options func(items *Items)

/*
ItemsPriceRange sets the exchange price range for request Items

https://api.dmarket.com/exchange/v1/market/items?priceFrom={priceFrom}&priceTo={priceTo}
*/
func ItemsPriceRange(priceFrom, priceTo int) Options {
	return func(i *Items) {
		if priceFrom < 0 || priceTo < priceFrom {
			panic(fmt.Errorf("%w [priceFrom %d priceTo %d] => priceFrom >= 0 && priceTo > priceFrom",
				ErrIncorrectPriceRange, priceFrom, priceTo))
		}
		i.priceFrom = priceFrom
		i.priceTo = priceTo
	}
}

/*
ItemsLimitPerRequest sets the exchange limit per request for Items

https://api.dmarket.com/exchange/v1/market/items?limit={limit}
*/
func ItemsLimitPerRequest(limit int) Options {
	return func(i *Items) {
		if limit <= 0 || limit > 100 {
			panic(ErrLimitPerRequest)
		}
		i.limit = limit
	}
}

/*
ItemsTitle sets the exchange limit per request for Items

https://api.dmarket.com/exchange/v1/market/items?title={title}
*/
func ItemsTitle(title string) Options {
	return func(i *Items) {
		i.title = title
	}
}

/*
GetAllItemsFromDmarket
gets all objects available on the Dmarket exchange with the parameters of the Client's Items and arguments.
The received response will be sent to the results channel.

If the answer from Dmarket contains an empty slice of schemas.Objects
(there are no objects on the market or inventory objects according to the specified filters have already been received),
but the HTTP response code is 200, will close the result channel.

If an error occurs during a request, parsing results, or receiving an HTTP error,
the error will be sent to the schemas.GetItemsResponse.Errors field, and the channel will be closed.

Available options:
	ItemsPriceRange(priceFrom, priceTo int)
	ItemsLimitPerRequest(limit int)
Panic when options get wrong options params!
*/
func (i Items) GetAllItemsFromDmarket(ctx context.Context, options ...Options) (results chan *GetItemsResponse) {
	return i.getAllItems(ctx, marketItems, options...)
}

/*
GetAllItemsFromUserInventory
gets all objects available on the user inventory with the parameters of the Items client and function arguments.
The received response will be sent to the results channel.

If the answer from Dmarket contains an empty slice of schemas.Objects
(there are no objects on the market or inventory objects according to the specified filters have already been received),
but the HTTP response code is 200, will close the result channel.

If an error occurs during a request, parsing results, or receiving an HTTP error,
the error will be sent to the schemas.GetItemsResponse.Errors field, and the channel will be closed.

Available options:
	ItemsPriceRange(priceFrom, priceTo int)
	ItemsLimitPerRequest(limit int)
Panic when options get wrong options params!
*/
func (i Items) GetAllItemsFromUserInventory(ctx context.Context, options ...Options) (results chan *GetItemsResponse) {
	return i.getAllItems(ctx, userItems, options...)
}

func (i *Items) getAllItems(ctx context.Context, from string, options ...Options) (results chan *GetItemsResponse) {
	for _, option := range options {
		option(i)
	}
	results = make(chan *GetItemsResponse, 1)
	go func() {
		defer close(results)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				items, err := i.GetItems(from)
				if err != nil || len(items.Objects) == 0 {
					results <- &GetItemsResponse{Error: err}
					return
				}
				results <- items
			}
		}
	}()
	return results
}

func (i *Items) GetItems(endpointURI string) (*GetItemsResponse, error) {
	params := &url.Values{
		"gameId":    {"9a92"},
		"currency":  {"USD"},
		"limit":     {strconv.Itoa(i.limit)},
		"priceFrom": {strconv.Itoa(i.priceFrom)},
		"priceTo":   {strconv.Itoa(i.priceTo)},
		"title":     {i.title},
		"cursor":    {i.cursor},
	}
	resp, err := i.client.Get(endpointURI + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("api (items): get items request error: %w", err)
	}
	var itemsResp GetItemsResponse
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api (items) error: %w http code: %s body: %s", ErrUnexpectedAPIResponse, resp.Status, resp.Body.String())
	}
	err = json.Unmarshal(resp.Body.Bytes(), &itemsResp)
	if err != nil {
		return nil, fmt.Errorf("api (items) error: %w into GetIntemsResponse struct "+
			"resp code: %s resp body: %s unmarshal error: %s", ErrUnmarshalAPIResponse, resp.Status, resp.Body.String(), err)
	}
	i.cursor = itemsResp.Cursor
	return &itemsResp, nil
}
