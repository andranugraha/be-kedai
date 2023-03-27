package repository

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/constant"
	errs "kedai/backend/be-kedai/internal/common/error"
	orderRepo "kedai/backend/be-kedai/internal/domain/order/repository"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShopRepository interface {
	FindShopById(id int) (*model.Shop, error)
	FindShopByUserId(userId int) (*model.Shop, error)
	FindShopBySlug(slug string) (*model.Shop, error)
	FindShopByKeyword(req dto.FindShopRequest) ([]*dto.FindShopResponse, int64, int, error)
	UpdateShopAddressIdByUserId(tx *gorm.DB, userId int, addressId int) error
	GetShopFinanceOverview(shopId int) (*dto.ShopFinanceOverviewResponse, error)
	GetShopStats(shopId int) (*dto.GetShopStatsResponse, error)
	GetShopInsight(shopId int, req dto.GetShopInsightRequest) (*dto.GetShopInsightResponse, error)
	FindShopByUserIdForUpdate(userId int) (*model.Shop, error)
	UpdateShop(shop *model.Shop) error
	Create(shop *model.Shop) error
	GetShopRating(userId int, filter dto.GetShopRatingFilterRequest) (*dto.ShopRatingResponse, int64, int, error)
}

type shopRepositoryImpl struct {
	db                 *gorm.DB
	invoicePerShopRepo orderRepo.InvoicePerShopRepository
	walletHistoryRepo  userRepo.WalletHistoryRepository
}

type ShopRConfig struct {
	DB                 *gorm.DB
	InvoicePerShopRepo orderRepo.InvoicePerShopRepository
	WalletHistoryRepo  userRepo.WalletHistoryRepository
}

func NewShopRepository(cfg *ShopRConfig) ShopRepository {
	return &shopRepositoryImpl{
		db:                 cfg.DB,
		invoicePerShopRepo: cfg.InvoicePerShopRepo,
		walletHistoryRepo:  cfg.WalletHistoryRepo,
	}
}

func (r *shopRepositoryImpl) FindShopById(id int) (*model.Shop, error) {
	var shop model.Shop

	err := r.db.Where("id = ?", id).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}

		return nil, err
	}

	return &shop, err
}

func (r *shopRepositoryImpl) FindShopByUserId(userId int) (*model.Shop, error) {
	var shop model.Shop

	err := r.db.Where("user_id = ?", userId).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}

		return nil, err
	}

	return &shop, err
}

func (r *shopRepositoryImpl) FindShopBySlug(slug string) (*model.Shop, error) {
	var shop model.Shop

	err := r.db.Where("slug = ?", slug).Preload("ShopCategory").First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}
	}

	return &shop, nil
}

func (r *shopRepositoryImpl) FindShopByKeyword(req dto.FindShopRequest) ([]*dto.FindShopResponse, int64, int, error) {
	var (
		shopList  []*dto.FindShopResponse
		totalRows int64
		totalPage int
		isActive  = true
	)

	db := r.db.Select(`shops.*, count(p.id) as product_count`).
		Joins("left join products p on shops.id = p.shop_id and p.is_active = ?", isActive).
		Group("shops.id").Where("shops.name ILIKE ?", "%"+req.Keyword+"%")

	countQuery := db.Session(&gorm.Session{})
	countQuery.Model(&model.Shop{}).Distinct("shops.id").Count(&totalRows)
	totalPage = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	err := db.Model(&model.Shop{}).Order("rating desc").Limit(req.Limit).Offset(req.Offset()).Find(&shopList).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return shopList, totalRows, totalPage, nil
}

func (r *shopRepositoryImpl) UpdateShopAddressIdByUserId(tx *gorm.DB, userId int, addressId int) error {
	res := tx.Model(&model.Shop{}).Where("user_id = ?", userId).Update("address_id", addressId)
	if err := res.Error; err != nil {
		tx.Rollback()
		return err
	}

	if res.RowsAffected == 0 {
		tx.Rollback()
		return errs.ErrShopNotFound
	}

	return nil
}

func (r *shopRepositoryImpl) GetShopFinanceOverview(shopId int) (*dto.ShopFinanceOverviewResponse, error) {
	toRelease, err := r.invoicePerShopRepo.GetShopFinanceToRelease(shopId)
	if err != nil {
		return nil, err
	}

	released, err := r.walletHistoryRepo.GetShopFinanceReleased(shopId)
	if err != nil {
		return nil, err
	}

	return &dto.ShopFinanceOverviewResponse{
		ToRelease: toRelease,
		Released:  *released,
	}, nil
}

func (r *shopRepositoryImpl) GetShopStats(shopId int) (*dto.GetShopStatsResponse, error) {
	var shopStats dto.GetShopStatsResponse

	err := r.db.Model(&model.Shop{}).
		Joins("left join invoice_per_shops ips on shops.id = ips.shop_id").
		Where("shops.id = ?", shopId).
		Select(`
				count(ips.id) filter (where ips.status = ?) as to_ship,
				count(ips.id) filter (where ips.status = ?) as shipping,
				count(ips.id) filter (where ips.status = ?) as completed,
				count(ips.id) filter (where ips.status = ?) as refund,
				(select count(p.id) from products p join skus s on p.id = s.product_id 
				where p.shop_id = shops.id and p.is_active = true
				group by p.id having sum(s.stock) = 0) as out_of_stock
		`, constant.TransactionStatusCreated,
			constant.TransactionStatusOnDelivery,
			constant.TransactionStatusCompleted,
			constant.TransactionStatusRefunded,
		).
		Group("shops.id").
		First(&shopStats).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}

		return nil, err
	}

	return &shopStats, nil
}

func (r *shopRepositoryImpl) GetShopInsight(shopId int, req dto.GetShopInsightRequest) (*dto.GetShopInsightResponse, error) {
	var (
		shopInsight dto.GetShopInsightResponse
	)

	row := r.db.Model(&model.Shop{}).
		Select(`(
				select count(sg.uuid) from shop_guests sg
				where sg.shop_id = shops.id
			) as visitor, 
			(
				select sum(p."view") from products p
				where p.shop_id = shops.id
			) as page_view, 
			(
				select count(ips.id) from invoice_per_shops ips
				where ips.shop_id = shops.id and ips.status = ?
			) as order
		`, constant.TransactionStatusCompleted).
		Where("shops.id = ?", shopId).
		Row()
	if row.Err() != nil {
		if errors.Is(row.Err(), gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}

		return nil, row.Err()
	}

	row.Scan(&shopInsight.Visitor, &shopInsight.PageView, &shopInsight.Order)

	sales, err := r.getShopSalesWithinInterval(shopId, req.Timeframe)
	if err != nil {
		return nil, err
	}

	shopInsight.Sales = sales

	return &shopInsight, nil
}

func (r *shopRepositoryImpl) getShopSalesWithinInterval(shopId int, timeframe string) ([]*dto.GetShopInsightSale, error) {
	var (
		shopInsightSale []*dto.GetShopInsightSale
		interval        string
		start           string
		end             string
		paymentDate     string
	)

	switch timeframe {
	case dto.ShopInsightTimeframeDay:
		interval = "2 hours"
		start = "date_trunc('day', now())"
		end = "date_trunc('day', now()) + interval '1 day'"
		paymentDate = "date_trunc('hour', i.payment_date)"
	case dto.ShopInsightTimeframeWeek:
		interval = "1 day"
		start = "date_trunc('week', now())"
		end = "date_trunc('week', now()) + interval '1 week'"
		paymentDate = "date_trunc('day', i.payment_date)"
	case dto.ShopInsightTimeframeMonth:
		interval = "1 week"
		start = "date_trunc('month', now())"
		end = "date_trunc('month', now()) + interval '1 month'"
		paymentDate = "date_trunc('week', i.payment_date)"
	}

	err := r.db.Model(&model.Shop{}).
		Raw(`
			WITH intervals AS (
				SELECT generate_series(
					`+start+` - interval '7 hours',
					`+end+` - interval '7 hours',
					interval '`+interval+`'
				) AS label
			)
			SELECT
				count(data.id) AS value,
				intervals.label
			FROM
				intervals
			LEFT JOIN (
				SELECT
					ips.id,
					`+paymentDate+` AS label
				FROM
					shops s
					LEFT JOIN invoice_per_shops ips ON s.id = ips.shop_id
					LEFT JOIN invoices i ON ips.invoice_id = i.id
				WHERE
					s.id = ? AND
					i.payment_date >= `+start+` - interval '7 hours'
			) AS data ON intervals.label <= data.label AND data.label < intervals.label + interval '`+interval+`'
			GROUP BY intervals.label
			ORDER BY intervals.label
		`, shopId).
		Scan(&shopInsightSale).Error
	if err != nil {
		return nil, err
	}

	return shopInsightSale, nil
}

func (r *shopRepositoryImpl) FindShopByUserIdForUpdate(userId int) (*model.Shop, error) {
	var shop model.Shop

	err := r.db.Model(&model.Shop{}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userId).
		First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}

		return nil, err
	}

	return &shop, nil
}

func (r *shopRepositoryImpl) UpdateShop(shop *model.Shop) error {
	return r.db.Save(shop).Error
}

func (r *shopRepositoryImpl) Create(shop *model.Shop) error {
	*shop.JoinedDate = time.Now()

	res := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(shop)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errs.ErrShopRegistered
	}

	return nil
}

func (r *shopRepositoryImpl) GetShopRating(userId int, filter dto.GetShopRatingFilterRequest) (*dto.ShopRatingResponse, int64, int, error) {
	whereClause := "s.user_id = ?"
	values := []interface{}{userId}

	if filter.Search != "" {
		whereClause += " AND p.name ILIKE ?"
		values = append(values, "%"+filter.Search+"%")
	}

	if filter.StartDate != "" && filter.EndDate != "" {
		whereClause += " AND tr.review_date BETWEEN ? AND ?"
		values = append(values, filter.StartDate, filter.EndDate)
	}

	if filter.Filter != 0 {
		whereClause += " AND tr.rating = ?"
		values = append(values, filter.Filter)
	}

	var results []*dto.ProductItem
	var ShopRating dto.ShopRatingResponse
	var totalRows int64
	var totalPages int

	query := r.db.Table("products p").
		Joins("JOIN product_medias pm ON p.id = pm.product_id").
		Joins("JOIN shops s ON s.id = p.shop_id").
		Joins("JOIN skus s2 ON s2.product_id = p.id").
		Joins("JOIN transactions t ON t.sku_id = s2.id").
		Joins("JOIN transaction_reviews tr ON tr.transaction_id = t.id").
		Joins("LEFT JOIN review_medias rm ON rm.review_id = tr.id").
		Joins("JOIN invoice_per_shops ips ON ips.id = t.invoice_id").
		Joins("JOIN users u ON u.id = t.user_id").
		Joins("JOIN user_profiles up ON up.user_id = u.id").
		Where(whereClause, values...).Distinct().
		Select(" t.id as transactionId, u.username as username, up.photo_url as user_url, ips.id as order_id, " +
			"(select url from product_medias pm2 where product_id =pm2.product_id limit 1) as product_url, " +
			"p.name as product_name, tr.rating as product_rating, tr.description as product_review, " +
			"string_agg(DISTINCT rm.url, ',') as review_medias, s.rating as shop_rating").
		Group("t.id, u.username, ips.id, up.photo_url, p.name, tr.rating, tr.description, s.rating ")

	r1 := query.Count(&totalRows)
	totalRows = r1.RowsAffected
	if r1.Error != nil {
		return nil, 0, 0, r1.Error
	}

	r2 := query.Limit(filter.Limit).Offset(filter.Limit * (filter.Page - 1)).Find(&results)
	if r2.Error != nil {
		return nil, 0, 0, r2.Error
	}

	for _, item := range results {
		if item.ReviewMedias != "" {
			item.ReviewMedia = strings.Split(item.ReviewMedias, ",")
			item.ReviewMedias = ""
		}
	}

	if len(results) == 0 {
		ShopRating.ShopRating = 0
	} else {
		ShopRating.ShopRating = results[0].ShopRating
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(filter.Limit)))
	ShopRating.ProductItem = results

	return &ShopRating, totalRows, totalPages, nil
}
