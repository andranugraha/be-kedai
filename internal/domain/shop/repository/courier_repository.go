package repository

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CourierRepository interface {
	GetAll() ([]*model.Courier, error)
	GetShipmentList(shopId int, request *dto.ShipmentCourierFilterRequest) ([]*dto.ShipmentCourierResponse, error)
	GetByShopID(shopID int) ([]*model.Courier, error)
	GetByServiceIDAndShopID(courierID, shopID int) (*model.Courier, error)
	GetByProductID(productID int) ([]*model.Courier, error)
	GetMatchingCouriersByShopIDAndProductIDs(*dto.MatchingProductCourierRequest) ([]*model.Courier, error)
	GetByID(id int) (*model.Courier, error)
	ToggleShopCourier(shopCouriers []*model.ShopCourier) error
	AddCourier(req *dto.ShipmentCourierRequest) (*model.Courier, error)
}

type courierRepositoryImpl struct {
	db                       *gorm.DB
	courierServiceRepository CourierServiceRepository
}

type CourierRConfig struct {
	DB                       *gorm.DB
	CourierServiceRepository CourierServiceRepository
}

func NewCourierRepository(cfg *CourierRConfig) CourierRepository {
	return &courierRepositoryImpl{
		db:                       cfg.DB,
		courierServiceRepository: cfg.CourierServiceRepository,
	}
}

func (r *courierRepositoryImpl) AddCourier(req *dto.ShipmentCourierRequest) (*model.Courier, error) {
	courier := &model.Courier{
		Name: req.Name,
		Code: req.Code,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(courier).Error
		if err != nil {
			return err
		}

		var serviceCouriers []*model.CourierService

		for _, service := range req.Service {
			serviceCouriers = append(serviceCouriers, &model.CourierService{
				CourierID:   courier.ID,
				Name:        service.Name,
				Code:        service.Code,
				MinDuration: service.MinDuration,
				MaxDuration: service.MaxDuration,
			})
		}

		err = r.courierServiceRepository.CreateCourierService(tx, serviceCouriers)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return courier, nil

}

func (r *courierRepositoryImpl) GetAll() ([]*model.Courier, error) {
	var couriers []*model.Courier

	err := r.db.Preload("Services").Find(&couriers).Error
	if err != nil {
		return nil, err
	}

	return couriers, nil
}

func (r *courierRepositoryImpl) GetShipmentList(shopId int, request *dto.ShipmentCourierFilterRequest) ([]*dto.ShipmentCourierResponse, error) {
	var couriers []*dto.ShipmentCourierResponse

	db := r.db.Select(`DISTINCT ON (couriers.id) couriers.*, COALESCE(sc.is_active, false) as is_active`).
		Joins("JOIN courier_services cs ON cs.courier_id = couriers.id").
		Joins("LEFT JOIN shop_couriers sc ON sc.courier_service_id = cs.id AND sc.shop_id = ?", shopId)

	if request.Status == constant.CourierStatusActive {
		db = db.Where("is_active")
	}
	if request.Status == constant.CourierStatusInactive {
		db = db.Not("is_active")
	}

	err := db.Model(&model.Courier{}).Find(&couriers).Error
	if err != nil {
		return nil, err
	}

	return couriers, nil
}

func (r *courierRepositoryImpl) GetByShopID(shopID int) ([]*model.Courier, error) {
	var couriers []*model.Courier

	err := r.db.
		Joins("JOIN courier_services cs ON couriers.id = cs.courier_id").
		Joins("JOIN shop_couriers ON cs.id = shop_couriers.courier_service_id").
		Where("shop_couriers.shop_id = ?", shopID).
		Distinct().
		Find(&couriers).
		Error

	if err != nil {
		return nil, err
	}

	return couriers, nil
}

func (r *courierRepositoryImpl) GetByServiceIDAndShopID(courierServiceID, shopID int) (*model.Courier, error) {
	var courier model.Courier

	err := r.db.
		Joins("JOIN courier_services cs ON couriers.id = cs.courier_id").
		Joins("JOIN shop_couriers ON cs.id = shop_couriers.courier_service_id").
		Where("shop_couriers.shop_id = ?", shopID).
		Where("cs.id = ?", courierServiceID).
		First(&courier).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErr.ErrCourierNotFound
		}

		return nil, err
	}

	return &courier, nil
}

func (r *courierRepositoryImpl) GetByProductID(productID int) ([]*model.Courier, error) {
	var couriers []*model.Courier

	err := r.db.
		Joins("JOIN courier_services cs ON couriers.id = cs.courier_id").
		Joins("JOIN product_couriers ON cs.id = product_couriers.courier_service_id").
		Where("product_couriers.product_id = ?", productID).
		Distinct().
		Find(&couriers).
		Error

	if err != nil {
		return nil, err
	}

	return couriers, nil
}

func (r *courierRepositoryImpl) GetMatchingCouriersByShopIDAndProductIDs(req *dto.MatchingProductCourierRequest) ([]*model.Courier, error) {
	var couriers []*model.Courier

	err := r.db.
		Preload("Services").
		Joins("JOIN courier_services cs ON couriers.id = cs.courier_id").
		Joins("JOIN shop_couriers sc on sc.courier_service_id = cs.id").
		Joins(`JOIN (SELECT courier_service_id
		FROM product_couriers
		WHERE product_id IN (?)
		GROUP BY courier_service_id
		HAVING COUNT(DISTINCT product_id) = ?)
		pc ON sc.courier_service_id = pc.courier_service_id`, req.ProductIDs, len(req.ProductIDs)).
		Where("sc.shop_id = ?", req.ShopID).
		Find(&couriers).
		Error

	if err != nil {
		return nil, err
	}

	return couriers, nil
}

func (r *courierRepositoryImpl) GetByID(id int) (*model.Courier, error) {
	var courier model.Courier

	err := r.db.Preload("Services").First(&courier, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErr.ErrCourierNotFound
		}

		return nil, err
	}

	return &courier, nil
}

func (r *courierRepositoryImpl) ToggleShopCourier(shopCouriers []*model.ShopCourier) error {
	return r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "courier_service_id"}, {Name: "shop_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"is_active": gorm.Expr("NOT shop_couriers.is_active")}),
		}).
		Clauses(clause.Returning{
			Columns: []clause.Column{{Name: "is_active"}},
		}).
		Create(&shopCouriers).Error
}
