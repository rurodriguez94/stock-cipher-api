package main

import (
	"time"

	"github.com/stock-cipher-api/encryption-api/internal/common"
	"github.com/stock-cipher-api/encryption-api/internal/security/handler"
	"github.com/stock-cipher-api/encryption-api/internal/security/service"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func main() {
	router := gin.Default()

	common.InitDefault()

	c := cache.New(5*time.Minute, 10*time.Minute)
	sh := handler.NewSecurityHandler(service.NewSecurityService(c))

	router.Use(common.ErrorHandler)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/security/encrypt", sh.EncryptData)
	router.POST("/security/decrypt/:token", sh.DecryptData)

	if err := router.Run(":8081"); err != nil {
		panic(err)
	}
}
