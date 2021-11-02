package items

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/defernest/dmarket-go/dmarket"
	"math/rand"
	"reflect"
	"strconv"
)

func init() {
	err := providers()
	if err != nil {
		panic(err)
	}
}

func (e *EndpointBehaviorOK) cursorValid(queryCursor string) bool {
	if queryCursor == e.cursor {
		e.cursor = faker.Password()
		return true
	}
	return false
}

func (q itemsQuery) generate(count int) []dmarket.Object {
	items := make([]dmarket.Object, count)
	for i := 0; i < len(items); i++ {
		err := faker.FakeData(&items[i])
		if err != nil {
			panic(err)
		}
		if q.Title != "" {
			items[i].Title = q.Title
			items[i].Extra.Name = q.Title
		}
		if q.PriceFrom == 0 && q.PriceTo == 0 {
			items[i].Price.Usd = strconv.Itoa(rand.Intn(10000000) + 1)
		} else {
			items[i].Price.Usd = strconv.Itoa(rand.Intn(q.PriceTo-q.PriceFrom+1) + q.PriceFrom)
		}
		items[i].GameID = q.GameId
		items[i].Extra.GameID = q.GameId
	}
	return items
}

func providers() error {
	err := faker.AddProvider("classID", func(v reflect.Value) (interface{}, error) {
		cID := rand.Intn(9999999999)
		iID := rand.Intn(9999999999)
		return fmt.Sprintf("%d:%d", iID, cID), nil
	})
	if err != nil {
		return err
	}
	err = faker.AddProvider("dprice", func(v reflect.Value) (interface{}, error) {
		price := strconv.Itoa(rand.Intn(99999))
		return price, nil
	})
	if err != nil {
		return err
	}
	err = faker.AddProvider("gems", func(v reflect.Value) (interface{}, error) {
		var gems = make([]dmarket.Gem, 2)
		for i := 0; i < len(gems); i++ {
			err := faker.FakeData(&gems[i])
			if err != nil {
				return nil, err
			}
		}
		return gems, err
	})
	if err != nil {
		return err
	}
	err = faker.AddProvider("stickers", func(v reflect.Value) (interface{}, error) {
		var stickers = make([]dmarket.Sticker, 2)
		for i := 0; i < len(stickers); i++ {
			err := faker.FakeData(&stickers[i])
			if err != nil {
				return nil, err
			}
		}
		return stickers, err
	})
	return err
}
