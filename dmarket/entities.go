package dmarket

type GetItemsResponse struct {
	Cursor  string   `json:"cursor"`
	Objects []Object `json:"objects"`
	Total   Total    `json:"total"`
	Error   error
}

type Total struct {
	Offers          int `json:"offers"`
	Targets         int `json:"targets"`
	Items           int `json:"items"`
	CompletedOffers int `json:"completedOffers"`
	ClosedTargets   int `json:"closedTargets"`
}

// Object represent entity.Object response from Dmarket
type Object struct {
	Amount             int64            `json:"amount" faker:"boundary_start=1, boundary_end=100"`
	ClassID            string           `json:"classId" faker:"classID"`
	CreatedAt          int64            `json:"createdAt" faker:"unix_time"`
	Description        string           `json:"description" faker:"paragraph"`
	Discount           int64            `json:"discount" faker:"boundary_start=0, boundary_end=99"`
	Extra              Extra            `json:"extra"`
	ExtraDoc           string           `json:"extraDoc" faker:"url"`
	GameID             string           `json:"gameId" faker:"-"`
	GameType           string           `json:"gameType" faker:"len=5"`
	Image              string           `json:"image" faker:"url"`
	InMarket           bool             `json:"inMarket"`
	Overpriced         int              `json:"overpriced" faker:"-"`
	InstantPrice       Price            `json:"instantPrice"`
	InstantTargetID    string           `json:"instantTargetId" faker:"uuid_hyphenated"`
	ItemID             string           `json:"itemId" faker:"uuid_hyphenated"`
	LockStatus         bool             `json:"lockStatus"`
	Owner              string           `json:"owner" faker:"uuid_hyphenated"`
	OwnerDetails       OwnerDetails     `json:"ownerDetails"`
	OwnersBlockchainID string           `json:"ownersBlockchainId" faker:"uuid_digit"`
	Price              Price            `json:"price"`
	RecommendedPrice   RecommendedPrice `json:"recommendedPrice"`
	Slug               string           `json:"slug" faker:"len=15"`
	Status             string           `json:"common" faker:"len=5"`
	SuggestedPrice     Price            `json:"suggestedPrice"`
	Title              string           `json:"title" faker:"len=25"`
	Type               string           `json:"type" faker:"len=5"`
}

type Extra struct {
	Ability           string    `json:"ability" faker:"len=10"`
	BackgroundColor   string    `json:"backgroundColor" faker:"len=5"`
	Category          string    `json:"category" faker:"len=10"`
	CategoryPath      string    `json:"categoryPath" faker:"len=25"`
	Class             []string  `json:"class" faker:"slice_len=4 len=10"`
	Collection        []string  `json:"collection" faker:"slice_len=4 len=10"`
	Exterior          string    `json:"exterior" faker:"len=10"`
	FloatValue        int64     `json:"floatValue"`
	GameID            string    `json:"gameId" faker:"-"`
	Gems              []Gem     `json:"gems" faker:"gems"`
	Grade             string    `json:"grade" faker:"len=10"`
	GroupID           string    `json:"groupId" faker:"uuid_hyphenated"`
	Growth            int64     `json:"growth" faker:"oneof: 0, 100"`
	Hero              string    `json:"hero" faker:"len=10"`
	InspectInGame     string    `json:"inspectInGame" faker:"url"`
	IsNew             bool      `json:"isNew"`
	ItemType          string    `json:"itemType" faker:"len=10"`
	LinkID            string    `json:"linkId" faker:"uuid_digit"`
	Name              string    `json:"name" faker:"len=10"`
	NameColor         string    `json:"nameColor" faker:"len=5"`
	OfferID           string    `json:"offerId" faker:"uuid_digit"`
	Quality           string    `json:"quality" faker:"len=5"`
	Rarity            string    `json:"rarity" faker:"len=5"`
	SerialNumber      int64     `json:"serialNumber" faker:"boundary_start=10000000, boundary_end=99999999"`
	Stickers          []Sticker `json:"stickers" faker:"stickers"`
	Subscribers       int64     `json:"subscribers" faker:"oneof: 0, 1000"`
	TagName           string    `json:"tagName" faker:"len=5"`
	Tradable          bool      `json:"tradable"`
	TradeLock         int64     `json:"tradeLock" faker:"oneof: 0, 100"`
	TradeLockDuration int64     `json:"tradeLockDuration" faker:"unix_time"`
	Type              string    `json:"type" faker:"len=5"`
	Videos            int64     `json:"videos" faker:"oneof: 0, 10"`
	ViewAtSteam       string    `json:"viewAtSteam" faker:"url"`
	Withdrawable      bool      `json:"withdrawable"`
}

type RecommendedPrice struct {
	D3     Price `json:"d3"`
	D7     Price `json:"d7"`
	D7Plus Price `json:"d7Plus"`
}

type OwnerDetails struct {
	Avatar string `json:"avatar"`
	ID     string `json:"id" faker:"uuid_hyphenated"`
	Wallet string `json:"wallet" faker:"uuid_digit"`
}

type Price struct {
	Dmc string `json:"DMC" faker:"dprice"`
	Usd string `json:"USD" faker:"dprice"`
}

type Gem struct {
	Image string `json:"image" faker:"url"`
	Name  string `json:"name" faker:"len=10"`
	Type  string `json:"type" faker:"len=10"`
}

type Sticker struct {
	Image string `json:"image" faker:"url"`
	Name  string `json:"name" faker:"len=10"`
}
