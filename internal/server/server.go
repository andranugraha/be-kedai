package server

import (
	"log"
	"time"

	"kedai/backend/be-kedai/connection"
	locationHandlerPackage "kedai/backend/be-kedai/internal/domain/location/handler"
	locationRepoPackage "kedai/backend/be-kedai/internal/domain/location/repository"
	locationServicePackage "kedai/backend/be-kedai/internal/domain/location/service"

	userRedisCache "kedai/backend/be-kedai/internal/domain/user/cache"
	userHandlerPackage "kedai/backend/be-kedai/internal/domain/user/handler"
	userRepoPackage "kedai/backend/be-kedai/internal/domain/user/repository"
	userServicePackage "kedai/backend/be-kedai/internal/domain/user/service"

	orderHandlerPackage "kedai/backend/be-kedai/internal/domain/order/handler"
	orderRepoPackage "kedai/backend/be-kedai/internal/domain/order/repository"
	orderServicePackage "kedai/backend/be-kedai/internal/domain/order/service"

	productHandlerPackage "kedai/backend/be-kedai/internal/domain/product/handler"
	productRepoPackage "kedai/backend/be-kedai/internal/domain/product/repository"
	productServicePackage "kedai/backend/be-kedai/internal/domain/product/service"

	shopHandlerPackage "kedai/backend/be-kedai/internal/domain/shop/handler"
	shopRepoPackage "kedai/backend/be-kedai/internal/domain/shop/repository"
	shopServicePackage "kedai/backend/be-kedai/internal/domain/shop/service"
	mail "kedai/backend/be-kedai/internal/utils/mail"
	random "kedai/backend/be-kedai/internal/utils/random"

	marketplaceHandlerPackage "kedai/backend/be-kedai/internal/domain/marketplace/handler"
	marketplaceRepoPackage "kedai/backend/be-kedai/internal/domain/marketplace/repository"
	marketplaceServicePackage "kedai/backend/be-kedai/internal/domain/marketplace/service"

	chatHandlerPackage "kedai/backend/be-kedai/internal/domain/chat/handler"
	chatRepoPackage "kedai/backend/be-kedai/internal/domain/chat/repository"
	chatServicePackage "kedai/backend/be-kedai/internal/domain/chat/service"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

func createRouter() *gin.Engine {
	db := connection.GetDB()
	redis := connection.GetCache()

	userCache := userRedisCache.NewUserCache(&userRedisCache.UserCConfig{
		RDC: redis,
	})
	walletCache := userRedisCache.NewWalletCache(&userRedisCache.WalletCConfig{
		RDC: redis,
	})

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
	maps := connection.GetGoogleMaps()

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

	walletHistoryRepo := userRepoPackage.NewWalletHistoryRepository(&userRepoPackage.WalletHistoryRConfig{
		DB: connection.GetDB(),
	})

	walletRepo := userRepoPackage.NewWalletRepository(&userRepoPackage.WalletRConfig{
		DB:            connection.GetDB(),
		WalletHistory: walletHistoryRepo,
	})

	courierRepo := shopRepoPackage.NewCourierRepository(&shopRepoPackage.CourierRConfig{
		DB: db,
	})

	invoiceStatusRepo := orderRepoPackage.NewInvoiceStatusRepository(&orderRepoPackage.InvoiceStatusRConfig{
		DB: db,
	})
	refundRequestRepo := orderRepoPackage.NewRefundRequestRepository(&orderRepoPackage.RefundRequestRConfig{
		DB: db,
	})

	skuRepo := productRepoPackage.NewSkuRepository(&productRepoPackage.SkuRConfig{
		DB: db})

	userCartItemRepo := userRepoPackage.NewUserCartItemRepository(&userRepoPackage.UserCartItemRConfig{
		DB: db,
	})

	invoiceRepo := orderRepoPackage.NewInvoiceRepository(&orderRepoPackage.InvoiceRConfig{
		DB:                db,
		UserCartItemRepo:  userCartItemRepo,
		SkuRepo:           skuRepo,
		UserWalletRepo:    walletRepo,
		InvoiceStatusRepo: invoiceStatusRepo,
		Redis:             userCache,
	})

	invoicePerShopRepo := orderRepoPackage.NewInvoicePerShopRepository(&orderRepoPackage.InvoicePerShopRConfig{
		DB:                db,
		WalletRepo:        walletRepo,
		InvoiceStatusRepo: invoiceStatusRepo,
		RefundRequestRepo: refundRequestRepo,
		SkuRepo:           skuRepo,
		UserVoucherRepo:   userVoucherRepo,
		InvoiceRepo:       invoiceRepo,
	})

	shopGuestRepo := shopRepoPackage.NewShopGuestRepository(&shopRepoPackage.ShopGuestRConfig{
		DB: db,
	})
	shopGuestService := shopServicePackage.NewShopGuestService(&shopServicePackage.ShopGuestSConfig{
		ShopGuestRepository: shopGuestRepo,
	})

	courierServiceRepo := shopRepoPackage.NewCourierServiceRepository(&shopRepoPackage.CourierServiceRConfig{
		DB: db,
	})
	courierServiceService := shopServicePackage.NewCourierServiceService(&shopServicePackage.CourierServiceSConfig{
		CourierServiceRepository: courierServiceRepo,
	})

	shopRepo := shopRepoPackage.NewShopRepository(&shopRepoPackage.ShopRConfig{
		DB:                 db,
		WalletHistoryRepo:  walletHistoryRepo,
		InvoicePerShopRepo: invoicePerShopRepo,
	})

	shopService := shopServicePackage.NewShopService(&shopServicePackage.ShopSConfig{
		ShopRepository:        shopRepo,
		CourierServiceService: courierServiceService,
	})

	courierService := shopServicePackage.NewCourierService(&shopServicePackage.CourierSConfig{
		CourierRepository: courierRepo,
		ShopService:       shopService,
	})

	shopVoucherRepo := shopRepoPackage.NewShopVoucherRepository(&shopRepoPackage.ShopVoucherRConfig{
		DB:                    db,
		UserVoucherRepository: userVoucherRepo,
	})

	shopVoucherService := shopServicePackage.NewShopVoucherService(&shopServicePackage.ShopVoucherSConfig{
		ShopVoucherRepository: shopVoucherRepo,
		ShopService:           shopService,
	})

	skuService := productServicePackage.NewSkuService(&productServicePackage.SkuSConfig{
		SkuRepository: skuRepo,
	})

	categoryRepo := productRepoPackage.NewCategoryRepository(&productRepoPackage.CategoryRConfig{
		DB: db,
	})

	categoryService := productServicePackage.NewCategoryService(&productServicePackage.CategorySConfig{
		CategoryRepo: categoryRepo,
	})

	variantGroupRepo := productRepoPackage.NewVariantGroupRepository(&productRepoPackage.VariantGroupRConfig{
		DB: db,
	})

	productVariantRepo := productRepoPackage.NewProductVariantRepository(&productRepoPackage.ProductVariantRConfig{
		DB: db,
	})

	productRepo := productRepoPackage.NewProductRepository(&productRepoPackage.ProductRConfig{
		DB:                       db,
		VariantGroupRepo:         variantGroupRepo,
		SkuRepository:            skuRepo,
		ProductVariantRepository: productVariantRepo,
	})
	productService := productServicePackage.NewProductService(&productServicePackage.ProductSConfig{
		ProductRepository:     productRepo,
		ShopVoucherService:    shopVoucherService,
		ShopService:           shopService,
		CourierService:        courierService,
		CategoryService:       categoryService,
		CourierServiceService: courierServiceService,
	})

	shopHandler := shopHandlerPackage.New(&shopHandlerPackage.HandlerConfig{
		ShopService:        shopService,
		ShopVoucherService: shopVoucherService,
		CourierService:     courierService,
		ShopGuestService:   shopGuestService,
	})

	refundRequestService := orderServicePackage.NewRefundRequestService(&orderServicePackage.RefundRequestSConfig{
		RefundRequestRepo: refundRequestRepo,
		ShopService:       shopService,
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

	walletService := userServicePackage.NewWalletService(&userServicePackage.WalletSConfig{
		WalletRepo:  walletRepo,
		UserCache:   userCache,
		WalletCache: walletCache,
		UserService: userService,
		MailUtils:   mailUtils,
		RandomUtils: randomUtils,
	})

	walletHistoryService := userServicePackage.NewWalletHistoryService(&userServicePackage.WalletHistorySConfig{
		WalletHistoryRepository: walletHistoryRepo,
		WalletService:           walletService,
	})

	userWishlistRepo := userRepoPackage.NewUserWishlistRepository(&userRepoPackage.UserWishlistRConfig{
		DB: db,
	})

	userWishlistService := userServicePackage.NewUserWishlistService(&userServicePackage.UserWishlistSConfig{
		UserWishlistRepository: userWishlistRepo,
		UserService:            userService,
		ProductService:         productService,
	})

	addressRepo := locationRepoPackage.NewAddressRepository(&locationRepoPackage.AddressRConfig{
		DB:              db,
		GoogleMaps:      maps,
		UserProfileRepo: userProfileRepo,
		ShopRepo:        shopRepo,
		SubdistrictRepo: subdistrictRepo,
	})

	addressService := locationServicePackage.NewAddressService(&locationServicePackage.AddressSConfig{
		AddressRepo:        addressRepo,
		ProvinceService:    provinceService,
		DistrictService:    districtService,
		SubdistrictService: subdistrictService,
		CityService:        cityService,
		UserProfileService: userProfileService,
		ShopService:        shopService,
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

	transactionReviewRepo := orderRepoPackage.NewTransactionReviewRepository(&orderRepoPackage.TransactionReviewRConfig{
		DB: db,
	})

	sealabsPayRepo := userRepoPackage.NewSealabsPayRepository(&userRepoPackage.SealabsPayRConfig{
		DB: db,
	})

	sealabsPayService := userServicePackage.NewSealabsPayService(&userServicePackage.SealabsPaySConfig{
		SealabsPayRepo: sealabsPayRepo,
	})

	userHandler := userHandlerPackage.New(&userHandlerPackage.HandlerConfig{
		UserService:          userService,
		WalletService:        walletService,
		WalletHistoryService: walletHistoryService,
		UserWishlistService:  userWishlistService,
		UserCartItemService:  userCartItemService,
		SealabsPayService:    sealabsPayService,
		AddressService:       addressService,
		UserProfileService:   userProfileService,
	})
	marketplaceHandler := marketplaceHandlerPackage.New(&marketplaceHandlerPackage.HandlerConfig{
		MarketplaceVoucherService: marketplaceVoucherService,
	})

	invoicePerShopService := orderServicePackage.NewInvoicePerShopService(&orderServicePackage.InvoicePerShopSConfig{
		InvoicePerShopRepo: invoicePerShopRepo,
		ShopService:        shopService,
		WalletService:      walletService,
	})
	transactionReviewService := orderServicePackage.NewTransactionReviewService(&orderServicePackage.TransactionReviewSConfig{
		TransactionReviewRepo: transactionReviewRepo,
		TransactionService:    transactionService,
		InvoicePerShopService: invoicePerShopService,
		ProductService:        productService,
	})
	productHandler := productHandlerPackage.New(&productHandlerPackage.Config{
		CategoryService:          categoryService,
		ProductService:           productService,
		SkuService:               skuService,
		TransactionReviewService: transactionReviewService,
	})

	invoiceService := orderServicePackage.NewInvoiceService(&orderServicePackage.InvoiceSConfig{
		InvoiceRepo:               invoiceRepo,
		AddressService:            addressService,
		ShopService:               shopService,
		ShopVoucherService:        shopVoucherService,
		CartItemService:           userCartItemService,
		ShopCourierService:        courierService,
		MarketplaceVoucherService: marketplaceVoucherService,
		SealabsPayService:         sealabsPayService,
		WalletService:             walletService,
	})

	orderHandler := orderHandlerPackage.New(&orderHandlerPackage.Config{
		InvoiceService:           invoiceService,
		TransactionReviewService: transactionReviewService,
		InvoicePerShopService:    invoicePerShopService,
		RefundRequestService:     refundRequestService,
	})

	locHandler := locationHandlerPackage.New(&locationHandlerPackage.Config{
		CityService:        cityService,
		ProvinceService:    provinceService,
		DistrictService:    districtService,
		SubdistrictService: subdistrictService,
		AddressService:     addressService,
	})

	chatHandler := chatHandlerPackage.New(&chatHandlerPackage.Config{
		ChatService: chatServicePackage.NewChatService(&chatServicePackage.ChatConfig{
			ChatRepo: chatRepoPackage.NewChatRepository(&chatRepoPackage.ChatRConfig{
				DB: db,
			}),
			ShopService: shopService,
			UserService: userService,
		}),
	})

	startCron(orderHandler)

	return NewRouter(&RouterConfig{
		UserHandler:        userHandler,
		LocationHandler:    locHandler,
		ProductHandler:     productHandler,
		ShopHandler:        shopHandler,
		MarketplaceHandler: marketplaceHandler,
		OrderHandler:       orderHandler,
		ChatHandler:        chatHandler,
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

func startCron(handler *orderHandlerPackage.Handler) {

	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(1).Hours().Do(func() {
		c := gin.Context{}

		handler.UpdateCronJob(&c)
	})

	if err != nil {
		log.Println(err)
	}

	scheduler.StartAsync()

}
