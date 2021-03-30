package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/stock-cipher-api/encryption-api/internal/apierror"
	"github.com/stock-cipher-api/encryption-api/internal/security"

	"github.com/gin-gonic/gin"
)

type TimeSeries struct {
	Date   time.Time `json:"date"`
	Open   string    `json:"open"`
	High   string    `json:"high"`
	Low    string    `json:"low"`
	Close  string    `json:"close"`
	Volume string    `json:"volume"`
}

type SecurityHandler struct {
	service security.Service
}

func NewSecurityHandler(service security.Service) *SecurityHandler {
	return &SecurityHandler{service: service}
}

func (h SecurityHandler) EncryptData(c *gin.Context) {
	var req interface{}

	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(apierror.NewStatusBadRequestError("cannot bind json"))
		return
	}

	res, err := h.service.EncryptData(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
	return
}

func (h SecurityHandler) DecryptData(c *gin.Context) {
	var req security.DecryptReq

	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(apierror.NewStatusBadRequestError("cannot bind json"))
		return
	}

	token := c.Param("token")
	if token == "" {
		_ = c.Error(apierror.NewStatusBadRequestError("token param cannot be empty"))
		return
	}

	res, err := h.service.DecryptData(c.Request.Context(), req.Payload, token)
	if err != nil {
		if errors.Is(err, security.ErrTokenNotFound) {
			_ = c.Error(apierror.NewStatusNotFoundError("token not found"))
			return
		}

		_ = c.Error(apierror.NewStatusUnauthorizedError("payload could not be decrypted"))
		return
	}

	c.JSON(http.StatusOK, res)
	return
}
