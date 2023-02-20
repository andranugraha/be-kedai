package dto

import (
	"gorm.io/gorm"
)

type GetCitiesRequest struct {
	Limit      int    `form:"limit"`
	Page       int    `form:"page"`
	ProvinceID int    `form:"provinceId"`
	Sort       string `form:"sort"`
}

func (r *GetCitiesRequest) Validate() {
	if r.Limit < 0 {
		r.Limit = 0
	}
	if r.Page < 1 {
		r.Page = 1
	}
}

func (r *GetCitiesRequest) Scope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ProvinceID != 0 {
			db = db.Where("cities.province_id = ?", r.ProvinceID)
		}
		if r.Sort == "most_shops" {
			db = db.Joins("left join user_addresses ua2 on cities.id = ua2.city_id and (select count(ua.id) from user_addresses ua inner join shops s on ua.id = s.address_id) > 0").
				Group("cities.id").
				Order("count(ua2.id) desc, cities.name asc")
		}
		return db
	}
}

func (r *GetCitiesRequest) Offset() int {
	return int((r.Page - 1) * r.Limit)
}
