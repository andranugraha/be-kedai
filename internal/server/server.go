package server

import (
	"log"
	"time"

	"kedai/backend/be-kedai/connection"
	locationRedisCache "kedai/backend/be-kedai/internal/domain/location/cache"
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

	productRedisCache "kedai/backend/be-kedai/internal/domain/product/cache"
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
	categoryCache := productRedisCache.NewCategoryCache(&productRedisCache.CategoryCConfig{
		RDC: redis,
	})
	locationCache := locationRedisCache.NewLocationCache(&locationRedisCache.LocationCConfig{
		RDC: redis,
	})

	userVoucherRepo := userRepoPackage.NewUserVoucherRepository(&userRepoPackage.UserVoucherRConfig{
		DB: db,
	})

	marketplaceVoucherRepo := marketplaceRepoPackage.NewMarketplaceVoucherRepository(&marketplaceRepoPackage.MarketplaceVoucherRConfig{
		DB:                    db,
		UserVoucherRepository: userVoucherRepo,
	})

	marketplaceBannerRepo := marketplaceRepoPackage.NewMarketplaceBannerRepository(&marketplaceRepoPackage.MarketplaceBannerRConfig{
		DB: db,
	})

	marketplaceVoucherService := marketplaceServicePackage.NewMarketplaceVoucherService(&marketplaceServicePackage.MarketplaceVoucherSConfig{
		MarketplaceVoucherRepository: marketplaceVoucherRepo,
	})

	marketplaceBannerService := marketplaceServicePackage.NewMarketplaceBannerService(&marketplaceServicePackage.MarketplaceBannerSConfig{
		MarketplaceBannerRepository: marketplaceBannerRepo,
	})

	mailer := connection.GetMailer()
	mailUtils := mail.NewMailUtils(&mail.MailUtilsConfig{Mailer: mailer})
	randomUtils := random.NewRandomUtils(&random.RandomUtilsConfig{})
	maps := connection.GetGoogleMaps()

	subdistrictRepo := locationRepoPackage.NewSubdistrictRepository(&locationRepoPackage.SubdistrictRConfig{
		DB: db,
	})
	subdistrictService := locationServicePackage.NewSubdistrictService(&locationServicePackage.SubdistrictSConfig{
		SubdistrictRepo: subdistrictRepo,
		Cache:           locationCache,
	})
	districtRepo := locationRepoPackage.NewDistrictRepository(&locationRepoPackage.DistrictRConfig{
		DB: db,
	})
	districtService := locationServicePackage.NewDistrictService(&locationServicePackage.DistrictSConfig{
		DistrictRepo: districtRepo,
		Cache:        locationCache,
	})
	cityRepo := locationRepoPackage.NewCityRepository(&locationRepoPackage.CityRConfig{
		DB: db,
	})
	cityService := locationServicePackage.NewCityService(&locationServicePackage.CitySConfig{
		CityRepo: cityRepo,
		Cache:    locationCache,
	})
	provinceRepo := locationRepoPackage.NewProvinceRepository(&locationRepoPackage.ProvinceRConfig{
		DB: db,
	})
	provinceService := locationServicePackage.NewProvinceService(&locationServicePackage.ProvinceSConfig{
		ProvinceRepo: provinceRepo,
		Cache:        locationCache,
	})

	walletHistoryRepo := userRepoPackage.NewWalletHistoryRepository(&userRepoPackage.WalletHistoryRConfig{
		DB: connection.GetDB(),
	})

	walletRepo := userRepoPackage.NewWalletRepository(&userRepoPackage.WalletRConfig{
		DB:            connection.GetDB(),
		WalletHistory: walletHistoryRepo,
	})

	invoiceStatusRepo := orderRepoPackage.NewInvoiceStatusRepository(&orderRepoPackage.InvoiceStatusRConfig{
		DB: db,
	})

	productPromotionRepo := productRepoPackage.NewProductPromotionRepository(&productRepoPackage.ProductPromotionRConfig{
		DB: db,
	})

	skuRepo := productRepoPackage.NewSkuRepository(&productRepoPackage.SkuRConfig{
		DB:                         db,
		ProductPromotionRepository: productPromotionRepo,
	})

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

	refundRequestRepo := orderRepoPackage.NewRefundRequestRepository(&orderRepoPackage.RefundRequestRConfig{
		DB: db,
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

	courierRepo := shopRepoPackage.NewCourierRepository(&shopRepoPackage.CourierRConfig{
		DB:                       db,
		CourierServiceRepository: courierServiceRepo,
	})

	shopVoucherRepo := shopRepoPackage.NewShopVoucherRepository(&shopRepoPackage.ShopVoucherRConfig{
		DB:                    db,
		UserVoucherRepository: userVoucherRepo,
	})

	skuService := productServicePackage.NewSkuService(&productServicePackage.SkuSConfig{
		SkuRepository: skuRepo,
	})

	categoryRepo := productRepoPackage.NewCategoryRepository(&productRepoPackage.CategoryRConfig{
		DB: db,
	})

	categoryService := productServicePackage.NewCategoryService(&productServicePackage.CategorySConfig{
		CategoryRepo:  categoryRepo,
		CategoryCache: categoryCache,
	})

	variantGroupRepo := productRepoPackage.NewVariantGroupRepository(&productRepoPackage.VariantGroupRConfig{
		DB: db,
	})

	productVariantRepo := productRepoPackage.NewProductVariantRepository(&productRepoPackage.ProductVariantRConfig{
		DB: db,
	})

	discussionRepo := productRepoPackage.NewDiscussionRepository(&productRepoPackage.DiscussionRConfig{
		DB: db,
	})

	productMediaRepo := productRepoPackage.NewProductMediaRepository(&productRepoPackage.ProductMediaRConfig{
		DB: db,
	})

	productRepo := productRepoPackage.NewProductRepository(&productRepoPackage.ProductRConfig{
		DB:                       db,
		VariantGroupRepo:         variantGroupRepo,
		SkuRepository:            skuRepo,
		ProductVariantRepository: productVariantRepo,
		DiscussionRepository:     discussionRepo,
		ProductMediaRepository:   productMediaRepo,
	})

	invoicePerShopRepo := orderRepoPackage.NewInvoicePerShopRepository(&orderRepoPackage.InvoicePerShopRConfig{
		DB:                db,
		WalletRepo:        walletRepo,
		InvoiceStatusRepo: invoiceStatusRepo,
		RefundRequestRepo: refundRequestRepo,
		SkuRepo:           skuRepo,
		ProductRepo:       productRepo,
		UserVoucherRepo:   userVoucherRepo,
		InvoiceRepo:       invoiceRepo,
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

	shopVoucherService := shopServicePackage.NewShopVoucherService(&shopServicePackage.ShopVoucherSConfig{
		ShopVoucherRepository: shopVoucherRepo,
		ShopService:           shopService,
	})

	courierService := shopServicePackage.NewCourierService(&shopServicePackage.CourierSConfig{
		CourierRepository: courierRepo,
		ShopService:       shopService,
	})

	discussionService := productServicePackage.NewDiscussionService(&productServicePackage.DiscussionSConfig{
		DiscussionRepository: discussionRepo,
		ShopService:          shopService,
	})

	productService := productServicePackage.NewProductService(&productServicePackage.ProductSConfig{
		ProductRepository:     productRepo,
		ShopVoucherService:    shopVoucherService,
		ShopService:           shopService,
		CourierService:        courierService,
		CategoryService:       categoryService,
		CourierServiceService: courierServiceService,
		DiscussionService:     discussionService,
	})

	shopPromotionRepo := shopRepoPackage.NewShopPromotionRepository(&shopRepoPackage.ShopPromotionRConfig{
		DB:                db,
		ProductRepository: productRepo,
	})

	shopPromotionService := shopServicePackage.NewShopPromotionService(&shopServicePackage.ShopPromotionSConfig{
		ShopPromotionRepository: shopPromotionRepo,
		ShopService:             shopService,
	})

	shopCategoryRepo := shopRepoPackage.NewShopCategoryRepository(&shopRepoPackage.ShopCategoryRConfig{
		DB:          db,
		ProductRepo: productRepo,
	})

	shopCategoryService := shopServicePackage.NewShopCategoryService(&shopServicePackage.ShopCategorySConfig{
		ShopCategoryRepo: shopCategoryRepo,
		ShopService:      shopService,
	})

	shopHandler := shopHandlerPackage.New(&shopHandlerPackage.HandlerConfig{
		ShopService:          shopService,
		ShopVoucherService:   shopVoucherService,
		ShopPromotionService: shopPromotionService,
		CourierService:       courierService,
		ShopGuestService:     shopGuestService,
		ShopCategoryService:  shopCategoryService,
	})

	userProfileRepo := userRepoPackage.NewUserProfileRepository(&userRepoPackage.UserProfileRConfig{
		DB: db,
	})

	userRepo := userRepoPackage.NewUserRepository(&userRepoPackage.UserRConfig{
		DB:              db,
		UserCache:       userCache,
		UserProfileRepo: userProfileRepo,
	})

	refundRequestRepo = orderRepoPackage.NewRefundRequestRepository(&orderRepoPackage.RefundRequestRConfig{
		DB:                 db,
		InvoicePerShopRepo: invoicePerShopRepo,
		InvoiceStatusRepo:  invoiceStatusRepo,
		UserRepo:           walletRepo,
		ProductRepo:        skuRepo,
	})

	refundRequestService := orderServicePackage.NewRefundRequestService(&orderServicePackage.RefundRequestSConfig{
		RefundRequestRepo: refundRequestRepo,
		ShopService:       shopService,
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
		MarketplaceBannerService:  marketplaceBannerService,
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
		DiscussionService:        discussionService,
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
			ShopService:    shopService,
			UserService:    userService,
			ProductService: productService,
			InvoiceService: invoicePerShopService,
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

	_, err = scheduler.Every(1).Minutes().Do(func() {
		c := gin.Context{}

		handler.ClearUnusedInvoice(&c)
	})

	if err != nil {
		log.Println(err)
	}

	scheduler.StartAsync()

}
