package handler

import (
	"errors"
	"net/http"

	"github.com/stock-cipher-api/stock-api/internal/apierror"
	"github.com/stock-cipher-api/stock-api/internal/stock"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	service stock.Service
}

func NewStockHandler(service stock.Service) *StockHandler {
	return &StockHandler{service: service}
}

func (h StockHandler) FetchData(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		_ = c.Error(apierror.NewStatusBadRequestError("symbol param cannot be empty"))
		return
	}

	res, err := h.service.GetStock(c.Request.Context(), symbol)
	if err != nil {
		if errors.Is(err, stock.ErrTimeSeriesNotFound) {
			_ = c.Error(apierror.NewStatusNotFoundError("time series daily not found"))
			return
		}
		if errors.Is(err, stock.ErrStockSymbolNotFound) {
			_ = c.Error(apierror.NewStatusNotFoundError("stock symbol not found"))
			return
		}
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h StockHandler) DecryptStockData(c *gin.Context) {
	token := c.Param("token")

	if token == "" {
		_ = c.Error(apierror.NewStatusBadRequestError("token param cannot be empty"))
		return
	}

	var req stock.DecryptStockParams
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(apierror.NewStatusBadRequestError("cannot bind json"))
		return
	}

	res, err := h.service.DecryptStockData(c.Request.Context(), req.Payload, token)
	if err != nil {
		if errors.Is(err, stock.ErrTokenNotFound) {
			_ = c.Error(apierror.NewStatusNotFoundError("token not found"))
			return
		}
		if errors.Is(err, stock.ErrCannotDecryptPayload) {
			_ = c.Error(apierror.NewStatusUnauthorizedError("unauthorized"))
			return
		}
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}
