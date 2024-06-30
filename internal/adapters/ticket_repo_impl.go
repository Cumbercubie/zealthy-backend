package repositories

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/zealthy/helpdesk-backend/internal/core/domain"
	"github.com/zealthy/helpdesk-backend/internal/core/ports"
	sql "github.com/zealthy/helpdesk-backend/internal/infra/database"
)

type ticketRepository struct {
	db *sql.DBHandler
}

func NewTicketRepository(db *sql.DBHandler) *ticketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ticketDetail ports.TicketCreateDetail) (*domain.Ticket, error) {
	// Implement database insertion logic
	err := validateTicketCreateDetail(ticketDetail)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	ticket, err := r.db.CreateTicket(ticketDetail)
	if err != nil {
		return nil, err
	}
	return ticket, err
}

func (r *ticketRepository) GetByID(id uuid.UUID) (*domain.Ticket, error) {
	ticket, err := r.db.GetTicketById(id)
	if err != nil {
		return nil, err
	}
	return ticket, err
}

func (r *ticketRepository) Update(ticketId uuid.UUID, options ports.TicketUpdateOptions) (*domain.Ticket, error) {
	// get current tickets
	currentTicket, err := r.db.GetTicketById(ticketId)

	if currentTicket == nil {
		return nil, fmt.Errorf("not found any ticket with this ID")
	}
	if err != nil {
		return nil, err
	}
	err = validateTicketUpdateOptions(currentTicket, options)

	if err != nil {
		return nil, fmt.Errorf("invalid update options: ", err.Error())
	}

	updatedTicket, err := r.db.UpdateTicket(ticketId, options)
	if err != nil {
		return nil, err
	}
	return updatedTicket, nil
}

func (r *ticketRepository) List(ticketOpts *ports.TicketListOptions) ([]*domain.Ticket, error) {
	tickets, err := r.db.ListTickets(ticketOpts)
	if err != nil {
		return nil, err
	}
	return tickets, err
}

func validateTicketUpdateOptions(ticket *domain.Ticket, opts ports.TicketUpdateOptions) error {
	// Allowed cases:
	// NEW -> IN_PROGRESS
	// NEW -> RESOLVED
	// IN_PROGRESS -> RESOLVED
	if opts.Status != nil && len(*opts.Status) >= 0 {
		if !opts.Status.IsValid() {
			return fmt.Errorf("invalid status")
		}
		switch *opts.Status {
		case domain.TicketStatusInProgress:
			if ticket.Status != domain.TicketStatusNew {
				return fmt.Errorf("invalid status: Can only move status from New to In Progress")
			}
		case domain.TicketStatusNew:
			return fmt.Errorf("invalid status: Ticket can not be moved back to New ")
		case domain.TicketStatusResolved:
			if ticket.Status == domain.TicketStatusResolved {
				return fmt.Errorf("invalid status: This ticket has already been resolved")
			}
		default:
			return nil
		}
	}
	return nil
}

func validateTicketCreateDetail(detail ports.TicketCreateDetail) error {
	if len(detail.IssuerEmail) <= 0 {
		return fmt.Errorf("invalid email")
	}

	if len(detail.Name) <= 0 {
		return fmt.Errorf("invalid name")
	}

	if len(*detail.Description) <= 0 {
		return fmt.Errorf("empty description")
	}
	return nil
}
