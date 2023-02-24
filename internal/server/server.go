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

	shopRepoPackage "kedai/backend/be-kedai/internal/domain/shop/repository"
	shopServicePackage "kedai/backend/be-kedai/internal/domain/shop/service"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	db := connection.GetDB()
	redis := connection.GetCache()

	productRepo := productRepoPackage.NewProductRepository(&productRepoPackage.ProductRConfig{
		DB: db,
	})
	productService := productServicePackage.NewProductService(&productServicePackage.ProductSConfig{
		ProductRepository: productRepo,
	})

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
		CityService: cityService,
	})

	walletRepo := userRepoPackage.NewWalletRepository(&userRepoPackage.WalletRConfig{
		DB: connection.GetDB(),
	})
	walletService := userServicePackage.NewWalletService(&userServicePackage.WalletSConfig{
		WalletRepo: walletRepo,
	})

	shopRepo := shopRepoPackage.NewShopRepository(&shopRepoPackage.ShopRConfig{
		DB: db,
	})

	shopService := shopServicePackage.NewShopService(&shopServicePackage.ShopSConfig{
		ShopRepository: shopRepo,
	})

	userRepo := userRepoPackage.NewUserRepository(&userRepoPackage.UserRConfig{
		DB: db,
	})

	userCache := userCache.NewUserCache(&userCache.UserCConfig{
		RDC: redis,
	})

	userService := userServicePackage.NewUserService(&userServicePackage.UserSConfig{
		Repository: userRepo,
		Redis:      userCache,
	})

	userWishlistRepo := userRepoPackage.NewUserWishlistRepository(&userRepoPackage.UserWishlistRConfig{
		DB: db,
	})

	userWishlistService := userServicePackage.NewUserWishlistService(&userServicePackage.UserWishlistSConfig{
		UserWishlistRepository: userWishlistRepo,
		UserService:            userService,
		ProductService:         productService,
	})
	skuRepo := productRepoPackage.NewSkuRepository(&productRepoPackage.SkuRConfig{
		DB: db,
	})

	skuService := productServicePackage.NewSkuService(&productServicePackage.SkuSConfig{
		SkuRepository: skuRepo,
	})

	userCartItemRepo := userRepoPackage.NewUserCartItemRepository(&userRepoPackage.UserCartItemRConfig{
		DB: db,
	})
	userAddressRepo := userRepoPackage.NewUserAddressRepository(&userRepoPackage.UserAddressRConfig{
		DB: db,
	})

	userAddressService := userServicePackage.NewUserAddressService(&userServicePackage.UserAddressSConfig{
		UserAddressRepo:    userAddressRepo,
		ProvinceService:    provinceService,
		DistrictService:    districtService,
		SubdistrictService: subdistrictService,
		CityService:        cityService,
	})

	userCartItemService := userServicePackage.NewUserCartItemService(&userServicePackage.UserCartItemSConfig{
		CartItemRepository: userCartItemRepo,
		SkuService:         skuService,
		ProductService:     productService,
		ShopService:        shopService,
	})

	userHandler := userHandlerPackage.New(&userHandlerPackage.HandlerConfig{
		UserService:         userService,
		WalletService:       walletService,
		UserWishlistService: userWishlistService,
		UserCartItemService: userCartItemService,
		UserAddressService:  userAddressService,
	})

	categoryRepo := productRepoPackage.NewCategoryRepository(&productRepoPackage.CategoryRConfig{
		DB: db,
	})

	categoryService := productServicePackage.NewCategoryService(&productServicePackage.CategorySConfig{
		CategoryRepo: categoryRepo,
	})

	productHandler := productHandlerPackage.New(&productHandlerPackage.Config{
		CategoryService: categoryService,
	})

	return NewRouter(&RouterConfig{
		UserHandler:     userHandler,
		LocationHandler: locHandler,
		ProductHandler:  productHandler,
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
