package handler

import (
	"kedai/backend/be-kedai/internal/common/code"
	spErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UserGetListOfChats(c *gin.Context) {
	userID := c.GetInt("userId")

	var param dto.ListOfChatsParamRequest
	c.ShouldBindQuery(&param)

	chatResponses, err := h.chatService.UserGetListOfChats(&param, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, spErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", chatResponses)
}

func (h *Handler) SellerGetListOfChats(c *gin.Context) {
	userID := c.GetInt("userId")

	var param dto.ListOfChatsParamRequest
	c.ShouldBindQuery(&param)

	chatResponses, err := h.chatService.SellerGetListOfChats(&param, userID)
	if err != nil {
		if err == spErr.ErrShopNotFound {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, spErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", chatResponses)
}

func (h *Handler) UserGetChat(c *gin.Context) {
	userID := c.GetInt("userId")

	shopSlug := c.Param("shopSlug")

	var param dto.ChatParamRequest
	c.ShouldBindQuery(&param)
	param.Validate()

	paginatedChats, err := h.chatService.UserGetChat(&param, userID, shopSlug)
	if err != nil {
		if err == spErr.ErrShopNotFound || err == spErr.ErrUserDoesNotExist {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, spErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", paginatedChats)
}

func (h *Handler) SellerGetChat(c *gin.Context) {
	userID := c.GetInt("userId")

	username := c.Param("username")

	var param dto.ChatParamRequest
	c.ShouldBindQuery(&param)
	param.Validate()

	paginatedChats, err := h.chatService.SellerGetChat(&param, userID, username)
	if err != nil {
		if err == spErr.ErrShopNotFound || err == spErr.ErrUserDoesNotExist {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, spErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", paginatedChats)
}

func (h *Handler) UserAddChat(c *gin.Context) {
	userID := c.GetInt("userId")

	shopSlug := c.Param("shopSlug")

	var body *dto.SendChatBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	chat, err := h.chatService.UserAddChat(body, userID, shopSlug)
	if err != nil {
		if err == spErr.ErrShopNotFound || err == spErr.ErrUserDoesNotExist || err == spErr.ErrSelfMessaging || err == spErr.ErrProductDoesNotExist || err == spErr.ErrInvoiceNotFound {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, spErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "success", chat)
}

func (h *Handler) SellerAddChat(c *gin.Context) {
	userID := c.GetInt("userId")

	username := c.Param("username")

	var body *dto.SendChatBodyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	chat, err := h.chatService.SellerAddChat(body, userID, username)
	if err != nil {
		if err == spErr.ErrShopNotFound || err == spErr.ErrUserDoesNotExist || err == spErr.ErrSelfMessaging || err == spErr.ErrProductDoesNotExist || err == spErr.ErrInvoiceNotFound {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, spErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "success", chat)
}
