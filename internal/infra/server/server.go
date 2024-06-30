package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/zealthy/helpdesk-backend/internal/core/ports"
)

type routeHandler struct {
	TicketService ports.TicketService
}

func NewServer() *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus((http.StatusOK))
			return
		}
		c.Next()
	})
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r
}
