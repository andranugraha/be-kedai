package server

import (
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/server/middleware"

	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"
	productHandler "kedai/backend/be-kedai/internal/domain/product/handler"
	shopHandler "kedai/backend/be-kedai/internal/domain/shop/handler"
	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	UserHandler     *userHandler.Handler
	LocationHandler *locationHandler.Handler
	ProductHandler  *productHandler.Handler
	ShopHandler     *shopHandler.Handler
}

func NewRouter(cfg *RouterConfig) *gin.Engine {
	r := gin.Default()

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = config.Origin
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsCfg.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsCfg.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(corsCfg))

	v1 := r.Group("/v1")
	{
		v1.Static("/docs", "swagger")

		user := v1.Group("/users")
		{
			user.POST("/register", cfg.UserHandler.UserRegistration)
			user.POST("/login", cfg.UserHandler.UserLogin)
			user.POST("/google-register", cfg.UserHandler.UserRegistrationWithGoogle)
			user.POST("/google-login", cfg.UserHandler.UserLoginWithGoogle)
			user.POST("/tokens/refresh", middleware.JWTValidateRefreshToken, cfg.UserHandler.RenewSession)
			userAuthenticated := user.Group("", middleware.JWTAuthorization, cfg.UserHandler.GetSession)
			{
				userAuthenticated.GET("", cfg.UserHandler.GetUserByID)
				userAuthenticated.POST("/logout", cfg.UserHandler.SignOut)

				userAuthenticated.PUT("/emails", cfg.UserHandler.UpdateUserEmail)
				userAuthenticated.PUT("/usernames", cfg.UserHandler.UpdateUsername)

				passwords := userAuthenticated.Group("/passwords")
				{
					passwords.POST("/change-request", cfg.UserHandler.RequestPasswordChange)
					passwords.POST("/change-confirmation", cfg.UserHandler.CompletePasswordChange)
				}

				profile := userAuthenticated.Group("/profiles")
				{
					profile.PUT("", cfg.UserHandler.UpdateProfile)
				}
				wallet := userAuthenticated.Group("/wallets")
				{
					wallet.GET("", cfg.UserHandler.GetWalletByUserID)
					wallet.POST("", cfg.UserHandler.RegisterWallet)
				}
				wishlists := userAuthenticated.Group("/wishlists")
				{
					wishlists.GET("", cfg.UserHandler.GetUserWishlists)
					wishlists.GET("/:productId", cfg.UserHandler.GetUserWishlist)
					wishlists.POST("", cfg.UserHandler.AddUserWishlist)
					wishlists.DELETE("/:productId", cfg.UserHandler.RemoveUserWishlist)
				}
				carts := userAuthenticated.Group("/carts")
				{
					carts.POST("", cfg.UserHandler.CreateCartItem)
					carts.GET("", cfg.UserHandler.GetAllCartItem)
					carts.PUT("/:skuId", cfg.UserHandler.UpdateCartItem)
				}
				addresses := userAuthenticated.Group("/addresses")
				{
					addresses.GET("", cfg.UserHandler.GetAllUserAddress)
					addresses.POST("", cfg.UserHandler.AddUserAddress)
					addresses.PUT("/:addressId", cfg.UserHandler.UpdateUserAddress)
					addresses.DELETE("/:addressId", cfg.UserHandler.DeleteUserAddress)
				}
				sealabsPay := userAuthenticated.Group("/sealabs-pays")
				{
					sealabsPay.GET("", cfg.UserHandler.GetSealabsPaysByUserID)
					sealabsPay.POST("", cfg.UserHandler.RegisterSealabsPay)
				}
			}
		}

		location := v1.Group("/locations")
		{
			location.GET("/cities", cfg.LocationHandler.GetCities)
			location.GET("/provinces", cfg.LocationHandler.GetProvinces)
			location.GET("/districts", cfg.LocationHandler.GetDistricts)
			location.GET("/subdistricts", cfg.LocationHandler.GetSubdistricts)
		}

		product := v1.Group("/products")
		{
			product.GET("/:code", cfg.ProductHandler.GetProductByCode)
			product.GET("/recommendations/categories", cfg.ProductHandler.GetRecommendationByCategory)
			category := product.Group("/categories")
			{
				category.GET("", cfg.ProductHandler.GetCategories)
			}
		}

		shop := v1.Group("/shops")
		{
			shop.GET("/:slug", cfg.ShopHandler.FindShopBySlug)
			shop.GET("/:slug/vouchers", cfg.ShopHandler.GetShopVoucher)
		}
	}

	return r
}
