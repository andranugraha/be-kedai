package handler

import (
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetInvoicePerShopsByUserID(c *gin.Context) {
	var request dto.InvoicePerShopFilterRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	request.Validate()

	userID := c.GetInt("userId")

	res, err := h.invoicePerShopService.GetInvoicesByUserID(userID, &request)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", res)
}
