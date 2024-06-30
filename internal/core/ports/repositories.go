package ports

import (
	"github.com/google/uuid"
	"github.com/zealthy/helpdesk-backend/internal/core/domain"
)

type TicketRepo interface {
	Create(ticketDetail TicketCreateDetail) (*domain.Ticket, error)
	GetByID(id uuid.UUID) (*domain.Ticket, error)
	Update(ticket uuid.UUID, opts TicketUpdateOptions) (*domain.Ticket, error)
	List(ticketOpts *TicketListOptions) ([]*domain.Ticket, error)
}
