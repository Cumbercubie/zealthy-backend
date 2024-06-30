package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zealthy/helpdesk-backend/internal/core/domain"
	"github.com/zealthy/helpdesk-backend/internal/core/ports"
)

type TicketHandler struct {
	service ports.TicketService
}

func NewTicketHandler(service ports.TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

func (rh *TicketHandler) RegisterTicketRoute(r *gin.Engine) {

	apiV1 := r.Group("/v1")

	apiV1.GET("/api/ticket/:id", rh.GetTicket)
	apiV1.POST("/api/ticket", rh.CreateTicket)
	apiV1.GET("/api/tickets", rh.ListTickets)
	apiV1.PUT("/api/ticket/:id", rh.UpdateTicket)

}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var input *domain.Ticket

	err := c.BindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.service.CreateTicket(ports.TicketCreateDetail{
		Name:        input.Name,
		IssuerEmail: input.IssuerEmail,
		Description: input.Description,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": []*domain.Ticket{ticket},
	})
}

func (h *TicketHandler) GetTicket(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.service.GetTicket(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Not found any ticket with this id",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": []*domain.Ticket{ticket},
	})
}

func (h *TicketHandler) ListTickets(c *gin.Context) {
	var params *ports.TicketListOptions = &ports.TicketListOptions{}
	err := c.BindQuery(&params)
	fmt.Println(c.Query("status"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	fmt.Println("params", params)
	tickets, err := h.service.ListTickets(params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(tickets) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Not found any tickets",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": tickets,
	})
}

func (h *TicketHandler) UpdateTicket(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input *ports.TicketUpdateOptions

	err = c.BindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.service.UpdateTicket(id, *input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ticket == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Could not retreive ticket",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": []*domain.Ticket{ticket},
	})
}
