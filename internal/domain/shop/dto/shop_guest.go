package dto

type AddShopGuestRequest struct {
	ShopId int `json:"shopId" binding:"required"`
}

func (req *AddShopGuestRequest) Validate() {
	if req.ShopId < 1 {
		req.ShopId = 0
	}
}
