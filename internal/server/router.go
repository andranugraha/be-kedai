package server

import (
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/server/middleware"

	locationHandler "kedai/backend/be-kedai/internal/domain/location/handler"
	marketplaceHandler "kedai/backend/be-kedai/internal/domain/marketplace/handler"
	orderHandler "kedai/backend/be-kedai/internal/domain/order/handler"
	productHandler "kedai/backend/be-kedai/internal/domain/product/handler"
	shopHandler "kedai/backend/be-kedai/internal/domain/shop/handler"
	userHandler "kedai/backend/be-kedai/internal/domain/user/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	UserHandler        *userHandler.Handler
	LocationHandler    *locationHandler.Handler
	ProductHandler     *productHandler.Handler
	ShopHandler        *shopHandler.Handler
	MarketplaceHandler *marketplaceHandler.Handler
	OrderHandler       *orderHandler.Handler
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

			user.POST("/passwords/reset-request", cfg.UserHandler.RequestPasswordReset)
			user.POST("/passwords/reset-confirmation", cfg.UserHandler.CompletePasswordReset)
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
					wallet.POST("/top-up", cfg.UserHandler.TopUp)
					wallet.GET("/histories/:ref", cfg.UserHandler.GetDetail)
					wallet.GET("/histories", cfg.UserHandler.GetWalletHistory)
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
					carts.DELETE("/:cartItemId", cfg.UserHandler.DeleteCartItem)
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
			product.GET("", cfg.ProductHandler.ProductSearchFiltering)
			product.GET("/:code", cfg.ProductHandler.GetProductByCode)
			product.GET("/:code/reviews", cfg.ProductHandler.GetProductReviews)
			product.GET("/:code/reviews/stats", cfg.ProductHandler.GetProductReviewStats)
			product.GET("/recommendations/categories", cfg.ProductHandler.GetRecommendationByCategory)
			category := product.Group("/categories")
			{
				category.GET("", cfg.ProductHandler.GetCategories)
			}
			sku := product.Group("/skus")
			{
				sku.GET("", cfg.ProductHandler.GetSKUByVariantIDs)
			}

		}

		shop := v1.Group("/shops")
		{
			shop.GET("", cfg.ShopHandler.FindShopByKeyword)
			shop.GET("/:slug", cfg.ShopHandler.FindShopBySlug)
			shop.GET("/:slug/products", cfg.ProductHandler.GetProductsByShopSlug)
			shop.GET("/:slug/vouchers", cfg.ShopHandler.GetShopVoucher)
			authenticated := shop.Group("", middleware.JWTAuthorization, cfg.UserHandler.GetSession)
			{
				authenticated.GET("/:slug/vouchers/valid", cfg.ShopHandler.GetValidShopVoucher)
			}
		}
		marketplace := v1.Group("/marketplaces")
		{
			marketplace.GET("/vouchers", cfg.MarketplaceHandler.GetMarketplaceVoucher)
			authenticated := marketplace.Group("", middleware.JWTAuthorization, cfg.UserHandler.GetSession)
			{
				authenticated.GET("/vouchers/valid", cfg.MarketplaceHandler.GetValidMarketplaceVoucher)
			}
		}

		order := v1.Group("/orders")
		{
			authenticated := order.Group("", middleware.JWTAuthorization, cfg.UserHandler.GetSession)
			{
				invoice := authenticated.Group("/invoices")
				{
					invoice.GET("", cfg.OrderHandler.GetInvoicePerShopsByUserID)
					invoice.GET("/:code", cfg.OrderHandler.GetInvoiceByCode)
				}

				transaction := authenticated.Group("/transactions")
				{
					transaction.GET("/:transactionId/reviews", cfg.OrderHandler.GetReviewByTransactionID)
					review := transaction.Group("/reviews")
					{
						review.POST("", cfg.OrderHandler.AddTransactionReview)
					}
				}
			}
		}
	}

	return r
}
