package ports

import (
	"github.com/google/uuid"
	"github.com/zealthy/helpdesk-backend/internal/core/domain"
)

type TicketCreateDetail struct {
	Name        string  `json:"name,omitempty"`
	IssuerEmail string  `json:"issuer_email,omitempty"`
	Description *string `json:"description,omitempty"`
	Response    *string `json:"response,omitempty"`
	Note        *string `json:"note,omitempty"`
}

type TicketQueryOptions struct {
	Offset string
	Limit  string
}

type TicketUpdateOptions struct {
	Status   *domain.TicketStatus `json:"status,omitempty"`
	Response *string              `json:"response,omitempty"`
	Note     *string              `json:"note,omitempty"`
}

type TicketListOptions struct {
	Status *domain.TicketStatus `form:"status"  json:"status,omitempty"`
}
type TicketService interface {
	CreateTicket(ticketDetail TicketCreateDetail) (*domain.Ticket, error)
	GetTicket(uuid uuid.UUID) (*domain.Ticket, error)
	UpdateTicket(id uuid.UUID, opts TicketUpdateOptions) (*domain.Ticket, error)
	ListTickets(ticketOpts *TicketListOptions) ([]*domain.Ticket, error)
}
