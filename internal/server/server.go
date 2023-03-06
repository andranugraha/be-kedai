package server

import (
	"log"

	"kedai/backend/be-kedai/connection"
	locationHandlerPackage "kedai/backend/be-kedai/internal/domain/location/handler"
	locationRepoPackage "kedai/backend/be-kedai/internal/domain/location/repository"
	locationServicePackage "kedai/backend/be-kedai/internal/domain/location/service"

	userCache "kedai/backend/be-kedai/internal/domain/user/cache"
	userHandlerPackage "kedai/backend/be-kedai/internal/domain/user/handler"
	userRepoPackage "kedai/backend/be-kedai/internal/domain/user/repository"
	userServicePackage "kedai/backend/be-kedai/internal/domain/user/service"

	productHandlerPackage "kedai/backend/be-kedai/internal/domain/product/handler"
	productRepoPackage "kedai/backend/be-kedai/internal/domain/product/repository"
	productServicePackage "kedai/backend/be-kedai/internal/domain/product/service"

	shopHandlerPackage "kedai/backend/be-kedai/internal/domain/shop/handler"
	shopRepoPackage "kedai/backend/be-kedai/internal/domain/shop/repository"
	shopServicePackage "kedai/backend/be-kedai/internal/domain/shop/service"

	orderHandlerPackage "kedai/backend/be-kedai/internal/domain/order/handler"
	orderRepoPackage "kedai/backend/be-kedai/internal/domain/order/repository"
	orderServicePackage "kedai/backend/be-kedai/internal/domain/order/service"
	mail "kedai/backend/be-kedai/internal/utils/mail"
	random "kedai/backend/be-kedai/internal/utils/random"

	marketplaceHandlerPackage "kedai/backend/be-kedai/internal/domain/marketplace/handler"
	marketplaceRepoPackage "kedai/backend/be-kedai/internal/domain/marketplace/repository"
	marketplaceServicePackage "kedai/backend/be-kedai/internal/domain/marketplace/service"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	db := connection.GetDB()
	redis := connection.GetCache()

	userVoucherRepo := userRepoPackage.NewUserVoucherRepository(&userRepoPackage.UserVoucherRConfig{
		DB: db,
	})

	marketplaceVoucherRepo := marketplaceRepoPackage.NewMarketplaceVoucherRepository(&marketplaceRepoPackage.MarketplaceVoucherRConfig{
		DB:                    db,
		UserVoucherRepository: userVoucherRepo,
	})

	marketplaceVoucherService := marketplaceServicePackage.NewMarketplaceVoucherService(&marketplaceServicePackage.MarketplaceVoucherSConfig{
		MarketplaceVoucherRepository: marketplaceVoucherRepo,
	})

	mailer := connection.GetMailer()
	mailUtils := mail.NewMailUtils(&mail.MailUtilsConfig{Mailer: mailer})
	randomUtils := random.NewRandomUtils(&random.RandomUtilsConfig{})

	districtRepo := locationRepoPackage.NewDistrictRepository(&locationRepoPackage.DistrictRConfig{
		DB: db,
	})
	districtService := locationServicePackage.NewDistrictService(&locationServicePackage.DistrictSConfig{
		DistrictRepo: districtRepo,
	})
	subdistrictRepo := locationRepoPackage.NewSubdistrictRepository(&locationRepoPackage.SubdistrictRConfig{
		DB: db,
	})
	subdistrictService := locationServicePackage.NewSubdistrictService(&locationServicePackage.SubdistrictSConfig{
		SubdistrictRepo: subdistrictRepo,
	})
	cityRepo := locationRepoPackage.NewCityRepository(&locationRepoPackage.CityRConfig{
		DB: db,
	})
	cityService := locationServicePackage.NewCityService(&locationServicePackage.CitySConfig{
		CityRepo: cityRepo,
	})
	provinceRepo := locationRepoPackage.NewProvinceRepository(&locationRepoPackage.ProvinceRConfig{
		DB: db,
	})
	provinceService := locationServicePackage.NewProvinceService(&locationServicePackage.ProvinceSConfig{
		ProvinceRepo: provinceRepo,
	})

	locHandler := locationHandlerPackage.New(&locationHandlerPackage.Config{
		CityService:        cityService,
		ProvinceService:    provinceService,
		DistrictService:    districtService,
		SubdistrictService: subdistrictService,
	})

	walletRepo := userRepoPackage.NewWalletRepository(&userRepoPackage.WalletRConfig{
		DB: connection.GetDB(),
	})
	walletService := userServicePackage.NewWalletService(&userServicePackage.WalletSConfig{
		WalletRepo: walletRepo,
	})

	courierRepo := shopRepoPackage.NewCourierRepository(&shopRepoPackage.CourierRConfig{
		DB: db,
	})
	courierService := shopServicePackage.NewCourierService(&shopServicePackage.CourierSConfig{
		CourierRepository: courierRepo,
	})

	shopRepo := shopRepoPackage.NewShopRepository(&shopRepoPackage.ShopRConfig{
		DB: db,
	})

	shopService := shopServicePackage.NewShopService(&shopServicePackage.ShopSConfig{
		ShopRepository: shopRepo,
	})

	shopVoucherRepo := shopRepoPackage.NewShopVoucherRepository(&shopRepoPackage.ShopVoucherRConfig{
		DB:                    db,
		UserVoucherRepository: userVoucherRepo,
	})

	shopVoucherService := shopServicePackage.NewShopVoucherService(&shopServicePackage.ShopVoucherSConfig{
		ShopVoucherRepository: shopVoucherRepo,
		ShopService:           shopService,
	})

	skuRepo := productRepoPackage.NewSkuRepository(&productRepoPackage.SkuRConfig{
		DB: db,
	})
	skuService := productServicePackage.NewSkuService(&productServicePackage.SkuSConfig{
		SkuRepository: skuRepo,
	})

	productRepo := productRepoPackage.NewProductRepository(&productRepoPackage.ProductRConfig{
		DB: db,
	})
	productService := productServicePackage.NewProductService(&productServicePackage.ProductSConfig{
		ProductRepository:  productRepo,
		ShopVoucherService: shopVoucherService,
		ShopService:        shopService,
		CourierService:     courierService,
	})

	shopHandler := shopHandlerPackage.New(&shopHandlerPackage.HandlerConfig{
		ShopService:        shopService,
		ShopVoucherService: shopVoucherService,
	})

	userCache := userCache.NewUserCache(&userCache.UserCConfig{
		RDC: redis,
	})
	userProfileRepo := userRepoPackage.NewUserProfileRepository(&userRepoPackage.UserProfileRConfig{
		DB: db,
	})

	userRepo := userRepoPackage.NewUserRepository(&userRepoPackage.UserRConfig{
		DB:              db,
		UserCache:       userCache,
		UserProfileRepo: userProfileRepo,
	})

	userService := userServicePackage.NewUserService(&userServicePackage.UserSConfig{
		Repository:  userRepo,
		Redis:       userCache,
		MailUtils:   mailUtils,
		RandomUtils: randomUtils,
	})

	userProfileService := userServicePackage.NewUserProfileService(&userServicePackage.UserProfileSConfig{
		Repository: userProfileRepo,
	})

	userWishlistRepo := userRepoPackage.NewUserWishlistRepository(&userRepoPackage.UserWishlistRConfig{
		DB: db,
	})

	userWishlistService := userServicePackage.NewUserWishlistService(&userServicePackage.UserWishlistSConfig{
		UserWishlistRepository: userWishlistRepo,
		UserService:            userService,
		ProductService:         productService,
	})

	userCartItemRepo := userRepoPackage.NewUserCartItemRepository(&userRepoPackage.UserCartItemRConfig{
		DB: db,
	})
	userAddressRepo := userRepoPackage.NewUserAddressRepository(&userRepoPackage.UserAddressRConfig{
		DB:              db,
		UserProfileRepo: userProfileRepo,
	})

	userAddressService := userServicePackage.NewUserAddressService(&userServicePackage.UserAddressSConfig{
		UserAddressRepo:    userAddressRepo,
		ProvinceService:    provinceService,
		DistrictService:    districtService,
		SubdistrictService: subdistrictService,
		CityService:        cityService,
		UserProfileService: userProfileService,
	})

	userCartItemService := userServicePackage.NewUserCartItemService(&userServicePackage.UserCartItemSConfig{
		CartItemRepository: userCartItemRepo,
		SkuService:         skuService,
		ProductService:     productService,
		ShopService:        shopService,
	})
	transactionRepo := orderRepoPackage.NewTransactionRepository(&orderRepoPackage.TransactionRConfig{
		DB: db,
	})

	transactionService := orderServicePackage.NewTransactionService(&orderServicePackage.TransactionSConfig{
		TransactionRepo: transactionRepo,
	})

	invoicePerShopRepo := orderRepoPackage.NewInvoicePerShopRepository(&orderRepoPackage.InvoicePerShopRConfig{
		DB: db,
	})

	invoicePerShopService := orderServicePackage.NewInvoicePerShopService(&orderServicePackage.InvoicePerShopSConfig{
		InvoicePerShopRepo: invoicePerShopRepo,
	})

	transactionReviewRepo := orderRepoPackage.NewTransactionReviewRepository(&orderRepoPackage.TransactionReviewRConfig{
		DB: db,
	})

	transactionReviewService := orderServicePackage.NewTransactionReviewService(&orderServicePackage.TransactionReviewSConfig{
		TransactionReviewRepo: transactionReviewRepo,
		TransactionService:    transactionService,
		InvoicePerShopService: invoicePerShopService,
	})

	sealabsPayRepo := userRepoPackage.NewSealabsPayRepository(&userRepoPackage.SealabsPayRConfig{
		DB: db,
	})

	sealabsPayService := userServicePackage.NewSealabsPayService(&userServicePackage.SealabsPaySConfig{
		SealabsPayRepo: sealabsPayRepo,
	})

	userHandler := userHandlerPackage.New(&userHandlerPackage.HandlerConfig{
		UserService:         userService,
		WalletService:       walletService,
		UserWishlistService: userWishlistService,
		UserCartItemService: userCartItemService,
		SealabsPayService:   sealabsPayService,
		UserAddressService:  userAddressService,
		UserProfileService:  userProfileService,
	})
	marketplaceHandler := marketplaceHandlerPackage.New(&marketplaceHandlerPackage.HandlerConfig{
		MarketplaceVoucherService: marketplaceVoucherService,
	})

	categoryRepo := productRepoPackage.NewCategoryRepository(&productRepoPackage.CategoryRConfig{
		DB: db,
	})

	categoryService := productServicePackage.NewCategoryService(&productServicePackage.CategorySConfig{
		CategoryRepo: categoryRepo,
	})

	productHandler := productHandlerPackage.New(&productHandlerPackage.Config{
		CategoryService: categoryService,
		ProductService:  productService,
		SkuService:      skuService,
	})

	orderHanlder := orderHandlerPackage.New(&orderHandlerPackage.Config{
		TransactionReviewService: transactionReviewService,
	})

	return NewRouter(&RouterConfig{
		UserHandler:        userHandler,
		LocationHandler:    locHandler,
		ProductHandler:     productHandler,
		ShopHandler:        shopHandler,
		MarketplaceHandler: marketplaceHandler,
		OrderHandler:       orderHanlder,
	})
}

func Init() {
	r := createRouter()

	err := r.Run()
	if err != nil {
		log.Println("error while running server", err)
		return
	}
}
