package service

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	productService "kedai/backend/be-kedai/internal/domain/product/service"
	shopService "kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserCartItemService interface {
	PreCheckCartItem(*dto.UserCartItemRequest) (*model.CartItem, *productModel.Sku, error)
	CreateCartItem(*dto.UserCartItemRequest) (*model.CartItem, error)
	GetAllCartItem(*dto.GetCartItemsRequest) (*commonDto.PaginationResponse, error)
	GetCartItemByIdAndUserId(int, int) (*model.CartItem, error)
}

type userCartItemServiceImpl struct {
	cartItemRepository repository.UserCartItemRepository
	skuService         productService.SkuService
	productService     productService.ProductService
	shopService        shopService.ShopService
}

type UserCartItemSConfig struct {
	CartItemRepository repository.UserCartItemRepository
	SkuService         productService.SkuService
	ProductService     productService.ProductService
	ShopService        shopService.ShopService
}

func NewUserCartItemService(cfg *UserCartItemSConfig) UserCartItemService {
	return &userCartItemServiceImpl{
		cartItemRepository: cfg.CartItemRepository,
		skuService:         cfg.SkuService,
		productService:     cfg.ProductService,
		shopService:        cfg.ShopService,
	}
}

func (s *userCartItemServiceImpl) CreateCartItem(cartItemReq *dto.UserCartItemRequest) (*model.CartItem, error) {

	result, sku, err := s.PreCheckCartItem(cartItemReq)

	// create new cart item if existing cart item not found
	if result == nil && sku != nil {
		cartItem := cartItemReq.ToUserCartItem()
		result, err = s.cartItemRepository.CreateCartItem(cartItem)
		if err != nil {
			return nil, err
		}

		return result, err
	}

	// update cart item quantity if existing cart item found
	if result != nil && sku != nil {
		result.Quantity += cartItemReq.Quantity
		if cartItemReq.Notes != "" {
			result.Notes = cartItemReq.Notes
		}

		result, err = s.cartItemRepository.UpdateCartItem(result)
		if err != nil {
			return nil, err
		}

		return result, err
	}

	return nil, err
}

func (s *userCartItemServiceImpl) PreCheckCartItem(cartItemReq *dto.UserCartItemRequest) (*model.CartItem, *productModel.Sku, error) {

	// Check sku (product) exist
	sku, err := s.skuService.GetByID(cartItemReq.SkuId)
	if err != nil {
		return nil, nil, err
	}

	// Check product exist
	product, err := s.productService.GetByID(sku.ProductId)
	if err != nil {
		return nil, nil, err
	}

	// check product active
	if !product.IsActive {
		return nil, nil, errs.ErrProductDoesNotExist
	}

	// check user not owner of shop
	shop, err := s.shopService.FindShopByUserId(cartItemReq.UserId)

	if err != nil && !errors.Is(err, errs.ErrShopNotFound) {
		return nil, nil, err
	}

	if shop != nil {
		if shop.UserID == cartItemReq.UserId {
			return nil, nil, errs.ErrUserIsShopOwner
		}
	}

	// Check if cart item already exists
	sameCartItem, err := s.cartItemRepository.GetCartItemByUserIdAndSkuId(cartItemReq.UserId, cartItemReq.SkuId)

	if err != nil && !errors.Is(err, errs.ErrCartItemNotFound) {
		return nil, nil, err
	}

	// update cart item quantity if existing cart item found
	if err == nil {
		if sameCartItem.Quantity+cartItemReq.Quantity > sku.Stock {
			return nil, nil, errs.ErrProductQuantityNotEnough
		}

		return sameCartItem, sku, nil
	}

	// return nil if existing cart item not found

	return nil, sku, nil
}

func (s *userCartItemServiceImpl) GetAllCartItem(cartItemReq *dto.GetCartItemsRequest) (*commonDto.PaginationResponse, error) {
	data := dto.GetCartItemsResponses{}

	cartItems, totalRows, totalPages, err := s.cartItemRepository.GetAllCartItem(cartItemReq)
	if err != nil {
		return nil, err
	}

	data.ToGetCartItemsResponses(cartItems)

	res := &commonDto.PaginationResponse{
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Limit:      cartItemReq.Limit,
		Page:       cartItemReq.Page,
		Data:       data.GetCartItemsResponses,
	}

	return res, nil
}

func (s *userCartItemServiceImpl) GetCartItemByIdAndUserId(id, userId int) (*model.CartItem, error) {
	return s.cartItemRepository.GetCartItemByIdAndUserId(id, userId)
}
