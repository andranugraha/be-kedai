package server

import (
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/connection"
	"kedai/backend/be-kedai/internal/server/middleware"

	"github.com/gin-contrib/pprof"

	chatHandler "kedai/backend/be-kedai/internal/domain/chat/handler"
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
	OrderHandler       *orderHandler.Handler
	MarketplaceHandler *marketplaceHandler.Handler
	ChatHandler        *chatHandler.Handler
}

func NewRouter(cfg *RouterConfig) *gin.Engine {
	r := gin.Default()

	pprof.Register(r)

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = config.Origin
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsCfg.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsCfg.ExposeHeaders = []string{"Content-Length"}
	corsCfg.AllowCredentials = true
	r.Use(cors.New(corsCfg))

	socketServer := connection.SocketIO()
	socket := r.Group("/socket.io")
	{
		socket.GET("/*any", gin.WrapH(socketServer))
		socket.POST("/*any", gin.WrapH(socketServer))
	}

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
					wallet.POST("/step-up", cfg.UserHandler.StepUp)
					wallet.POST("/pins/change-requests", cfg.UserHandler.RequestWalletPinChange)
					wallet.POST("/pins/change-confirmations", cfg.UserHandler.CompleteChangeWalletPin)
					wallet.POST("/pins/reset-requests", cfg.UserHandler.RequestWalletPinReset)
					wallet.POST("/pins/reset-confirmations", cfg.UserHandler.CompleteResetWalletPin)
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
					carts.DELETE("", cfg.UserHandler.DeleteCartItem)
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
				chat := userAuthenticated.Group("/chats")
				{
					chat.GET("/", cfg.ChatHandler.UserGetListOfChats)
					chat.GET("/:shopSlug", cfg.ChatHandler.UserGetChat)
					chat.POST("/:shopSlug", cfg.ChatHandler.UserAddChat)
				}
			}
		}

		location := v1.Group("/locations")
		{
			location.GET("/cities", cfg.LocationHandler.GetCities)
			location.GET("/provinces", cfg.LocationHandler.GetProvinces)
			location.GET("/districts", cfg.LocationHandler.GetDistricts)
			location.GET("/subdistricts", cfg.LocationHandler.GetSubdistricts)
			location.GET("/addresses", cfg.LocationHandler.SearchAddress)
			location.GET("/addresses/:placeId", cfg.LocationHandler.SearchAddressDetail)
		}

		product := v1.Group("/products")
		{
			product.GET("", cfg.ProductHandler.ProductSearchFiltering)
			product.GET("/:code", cfg.ProductHandler.GetProductByCode)
			product.GET("/:code/reviews", cfg.ProductHandler.GetProductReviews)
			product.GET("/:code/reviews/stats", cfg.ProductHandler.GetProductReviewStats)
			product.GET("/recommendations/categories", cfg.ProductHandler.GetRecommendationByCategory)
			product.GET("/autocompletes", cfg.ProductHandler.SearchAutocomplete)
			product.GET("/recommended", cfg.ProductHandler.GetRecommendedProducts)
			product.POST("/views", cfg.ProductHandler.AddProductView)
			product.GET("/discussions/:productId", cfg.ProductHandler.GetDiscussionByProductID)
			product.GET("/discussions/replies/:parentId", cfg.ProductHandler.GetDiscussionByParentID)
			product.POST("/discussions", middleware.JWTAuthorization, cfg.ProductHandler.PostDiscussion)
			category := product.Group("/categories")
			{
				category.GET("", cfg.ProductHandler.GetCategories)
			}
			sku := product.Group("/skus")
			{
				sku.GET("", cfg.ProductHandler.GetSKUByVariantIDs)
			}

			authenticated := product.Group("", middleware.JWTAuthorization, cfg.UserHandler.GetSession)
			{
				authenticated.POST("", cfg.ProductHandler.CreateProduct)
			}

		}

		shop := v1.Group("/shops")
		{
			shop.GET("", cfg.ShopHandler.FindShopByKeyword)
			visitor := shop.Group("/visitors")
			{
				visitor.POST("", cfg.ShopHandler.AddShopGuest)
			}
			shop.GET("/:slug", cfg.ShopHandler.FindShopBySlug)
			shop.GET("/:slug/products", cfg.ProductHandler.GetProductsByShopSlug)
			shop.GET("/:slug/vouchers", cfg.ShopHandler.GetShopVoucher)
			authenticated := shop.Group("", middleware.JWTAuthorization, cfg.UserHandler.GetSession)
			{
				authenticated.GET("/profile", cfg.ShopHandler.GetShopProfile)
				authenticated.PUT("/profile", cfg.ShopHandler.UpdateShopProfile)
				authenticated.GET("/:slug/vouchers/valid", cfg.ShopHandler.GetValidShopVoucher)
				authenticated.GET("/:slug/couriers", cfg.ShopHandler.GetMatchingCouriers)
			}
		}
		marketplace := v1.Group("/marketplaces")
		{
			marketplace.GET("/vouchers", cfg.MarketplaceHandler.GetMarketplaceVoucher)
			marketplace.POST("/couriers", cfg.ShopHandler.AddCourier)
			marketplace.GET("/banners", cfg.MarketplaceHandler.GetMarketplaceBanner)
			authenticated := marketplace.Group("", middleware.JWTAuthorization, cfg.UserHandler.GetSession)
			{
				authenticated.GET("/vouchers/valid", cfg.MarketplaceHandler.GetValidMarketplaceVoucher)
				authenticated.GET("/couriers", cfg.ShopHandler.GetAllCouriers)
			}
		}

		order := v1.Group("/orders")
		{
			authenticated := order.Group("", middleware.JWTAuthorization, cfg.UserHandler.GetSession)
			{
				authenticated.POST("", cfg.OrderHandler.Checkout)
				invoice := authenticated.Group("/invoices")
				{
					invoice.POST("", cfg.OrderHandler.PayInvoice)
					invoice.POST("/cancel", cfg.OrderHandler.CancelCheckout)
					invoice.GET("", cfg.OrderHandler.GetInvoicePerShopsByUserID)
					invoice.GET("/:code", cfg.OrderHandler.GetInvoiceByCode)
					invoice.PUT("/:code/receive", cfg.OrderHandler.UpdateToReceived)
					invoice.PUT("/:code/complete", cfg.OrderHandler.UpdateToCompleted)
					invoice.POST("/:code/refund", cfg.OrderHandler.Refund)
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

		// TODO ADD MIDLEWARE FOR AUTH ADMIN
		admin := v1.Group("/admins")
		{
			admin.POST("/login", cfg.UserHandler.AdminSignIn)
			authenticated := admin.Group("", middleware.AdminJWTAuthorization, cfg.UserHandler.GetSession)
			{
				category := authenticated.Group("/categories")
				{
					category.POST("", cfg.ProductHandler.AddCategory)
				}
				order := authenticated.Group("/orders")
				{
					order.POST("/:orderId/cancel-commit", cfg.OrderHandler.UpdateToCanceled)
				}
				marketplace := authenticated.Group("/marketplaces")
				{
					marketplace.POST("/banners", cfg.MarketplaceHandler.AddMarketplaceBanner)
					voucher := marketplace.Group("/vouchers")
					{
						voucher.POST("", cfg.MarketplaceHandler.CreateMarketplaceVoucher)
						voucher.GET("", cfg.MarketplaceHandler.GetMarketplaceVoucherAdmin)
						voucher.GET("/:code", cfg.MarketplaceHandler.GetMarketplaceVoucherAdminByCode)
					}
				}
			}
		}

		seller := v1.Group("/sellers")
		{
			authenticated := seller.Group("", middleware.JWTAuthorization, cfg.UserHandler.GetSession)
			{
				authenticated.POST("/register", cfg.ShopHandler.CreateShop)
				authenticated.GET("/stats", cfg.ShopHandler.GetShopStats)
				authenticated.GET("/insights", cfg.ShopHandler.GetShopInsights)
				authenticated.GET("/ratings", cfg.ShopHandler.GetShopRating)
				authenticated.GET("/discussions", cfg.ProductHandler.GetUnrepliedDiscussionByShopID)
				finance := authenticated.Group("/finances")
				{
					income := finance.Group("/incomes")
					{
						income.GET("", cfg.OrderHandler.GetInvoicePerShopsByShopId)
						income.GET("/overviews", cfg.ShopHandler.GetShopFinanceOverview)
						income.POST("/withdrawals", cfg.OrderHandler.WithdrawFromInvoice)
						income.GET("/:orderId", cfg.OrderHandler.GetInvoiceByShopIdAndOrderId)
					}
				}
				courier := authenticated.Group("/couriers")
				{
					courier.GET("", cfg.ShopHandler.GetShipmentList)
					courier.POST("", cfg.ShopHandler.ToggleShopCourier)
				}

				product := authenticated.Group("/products")
				{
					product.GET("", cfg.ProductHandler.GetSellerProducts)
					product.GET("/:code", cfg.ProductHandler.GetSellerProductDetailByCode)
					product.PUT("/:code", cfg.ProductHandler.UpdateProduct)
					product.PUT("/:code/activations", cfg.ProductHandler.UpdateProductActivation)
				}

				voucher := authenticated.Group("/vouchers")
				{
					voucher.GET("", cfg.ShopHandler.GetSellerVoucher)
					voucher.GET("/:code", cfg.ShopHandler.GetVoucherByCodeAndShopId)
					voucher.POST("", cfg.ShopHandler.CreateVoucher)
					voucher.PUT("/:code", cfg.ShopHandler.UpdateVoucher)
					voucher.DELETE("/:code", cfg.ShopHandler.DeleteVoucher)
				}

				promotion := authenticated.Group("/promotions")
				{
					promotion.GET("", cfg.ShopHandler.GetSellerPromotions)
					promotion.GET("/:promotionId", cfg.ShopHandler.GetSellerPromotionById)
					promotion.PUT("/:promotionId", cfg.ShopHandler.UpdatePromotion)
					promotion.POST("", cfg.ShopHandler.CreateShopPromotion)
					promotion.DELETE("/:promotionId", cfg.ShopHandler.DeletePromotion)
				}

				order := authenticated.Group("/orders")
				{
					order.GET("", cfg.OrderHandler.GetShopOrder)
					order.GET("/:orderId", cfg.OrderHandler.GetInvoiceByShopIdAndOrderId)
					order.PUT("/:orderId/process", cfg.OrderHandler.UpdateToProcessing)
					order.PUT("/:orderId/delivery", cfg.OrderHandler.UpdateToDelivery)
					order.POST("/:orderId/cancel-request", cfg.OrderHandler.UpdateToRefundPendingSellerCancel)
					order.PUT("/:orderId/refund", cfg.OrderHandler.UpdateRefundStatus)
				}
				chat := authenticated.Group("/chats")
				{
					chat.GET("/", cfg.ChatHandler.SellerGetListOfChats)
					chat.GET("/:username", cfg.ChatHandler.SellerGetChat)
					chat.POST("/:username", cfg.ChatHandler.SellerAddChat)
				}
				category := authenticated.Group("/categories")
				{
					category.GET("", cfg.ShopHandler.GetSellerCategories)
				}
			}
		}
	}

	return r
}
