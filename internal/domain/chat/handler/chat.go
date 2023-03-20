package handler

import (
	"kedai/backend/be-kedai/internal/common/code"
	spErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddChat(c *gin.Context) {
	userID := c.GetInt("userId")

	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	chat, err := h.chatService.AddChat(req.RoomId, req.Message, req.Type, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, spErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", chat)
}
