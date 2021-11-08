package items

import (
	"github.com/defernest/dmarket-go/dmarket"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Params struct {
	GameId    string `form:"gameId" binding:"required"`
	Title     string `form:"title"`
	Currency  string `form:"currency" binding:"required,contains=USD"`
	Cursor    string `form:"cursor"`
	Limit     int    `form:"limit" binding:"required,gte=0,lte=100"`
	PriceFrom int    `form:"priceFrom" binding:"gte=0"`
	PriceTo   int    `form:"priceTo" binding:"gtefield=PriceFrom"`
}

type EndpointBehaviorOK struct {
	count  int
	cursor string
}

func (e *EndpointBehaviorOK) Endpoint() (httpMethod string, relativePath string, handler gin.HandlerFunc) {
	return http.MethodGet, "/exchange/v1/market/items", func(context *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				context.String(dmarket.ErrorRepresentation{Response: dmarket.Response{StatusCode: http.StatusInternalServerError}}.String())
			}
		}()

		var itemsQuery Params
		err := context.ShouldBindQuery(&itemsQuery)
		if err != nil || !e.cursorValid(itemsQuery.Cursor) {
			context.String(dmarket.ErrorRepresentation{Response: dmarket.Response{StatusCode: http.StatusBadRequest}}.String())
			return
		}

		resp := dmarket.GetItemsResponse{Total: dmarket.Total{Items: e.count}, Cursor: e.cursor}
		if (e.count - itemsQuery.Limit) >= 0 {
			resp.Objects = itemsQuery.GenerateItems(itemsQuery.Limit)
		} else {
			resp.Objects = itemsQuery.GenerateItems(e.count)
		}
		context.JSON(http.StatusOK, &resp)
		e.count -= itemsQuery.Limit
		if e.count < 0 {
			e.count = 0
		}
	}
}

func MustReturnSuccess(count int) *EndpointBehaviorOK {
	return &EndpointBehaviorOK{count: count}
}
