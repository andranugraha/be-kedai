package dto

type FindShopRequest struct {
	Keyword string `form:"keyword"`
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
}

type FindShopResponse struct {
	Slug         string  `json:"slug"`
	Name         string  `json:"name"`
	ProductCount int64   `json:"productCount"`
	Rating       float64 `json:"rating"`
	PhotoUrl     string  `json:"photoUrl"`
}

type ShopFinanceOverviewResponse struct {
	ToRelease float64             `json:"toRelease"`
	Released  ShopFinanceReleased `json:"released"`
}

type ShopFinanceReleased struct {
	Week  float64 `json:"week"`
	Month float64 `json:"month"`
	Total float64 `json:"total"`
}

func (req *FindShopRequest) Validate() {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 10
	}
}

func (req *FindShopRequest) Offset() int {
	return (req.Page - 1) * req.Limit
}

type GetShopStatsResponse struct {
	ToShip     int `json:"toShip"`
	Shipping   int `json:"shipping"`
	Completed  int `json:"completed"`
	Refund     int `json:"refund"`
	OutOfStock int `json:"outOfStock"`
}

type GetShopInsightRequest struct {
	Timeframe string `form:"timeframe"`
	UserId    int
}

type GetShopInsightResponse struct {
	Visitor  int                   `json:"visitor"`
	PageView int                   `json:"pageView"`
	Order    int                   `json:"order"`
	Sales    []*GetShopInsightSale `json:"sales"`
}

type GetShopInsightSale struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

const (
	ShopInsightTimeframeDay   = "day"
	ShopInsightTimeframeWeek  = "week"
	ShopInsightTimeframeMonth = "month"
)

func (req *GetShopInsightRequest) Validate() {
	if req.Timeframe == "" {
		req.Timeframe = ShopInsightTimeframeDay
	}
}
