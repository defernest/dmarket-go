package dmarket

type Exchange struct {
	Items *Items
}

/*
NewExchange create new Exchange endpoint client with default params
	Items{
		client:        client,
		priceFrom:     0,
		priceTo:       1000000,
		limit:         100,
	}
*/
func NewExchange(client Requester) *Exchange {
	exchange := &Exchange{
		Items: &Items{
			client:    client,
			priceFrom: 0,
			priceTo:   1000000,
			limit:     100,
		}}
	return exchange
}
