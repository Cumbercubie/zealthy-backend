package services

import (
	"github.com/google/uuid"
	"github.com/zealthy/helpdesk-backend/internal/core/domain"
	"github.com/zealthy/helpdesk-backend/internal/core/ports"
)

type ticketService struct {
	repo ports.TicketRepo
}
type ListTicketOpts struct {
	limit  int64
	offset int64
}

func NewTicketService(repo ports.TicketRepo) ports.TicketService {
	return &ticketService{repo: repo}
}

func (s *ticketService) CreateTicket(ticketDetail ports.TicketCreateDetail) (*domain.Ticket, error) {
	result, err := s.repo.Create(ports.TicketCreateDetail{
		Name:        ticketDetail.Name,
		IssuerEmail: ticketDetail.IssuerEmail,
		Description: ticketDetail.Description,
		Response:    nil,
		Note:        nil,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ticketService) GetTicket(uuid uuid.UUID) (*domain.Ticket, error) {
	ticket, err := s.repo.GetByID(uuid)

	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *ticketService) ListTickets(ticketOpts *ports.TicketListOptions) ([]*domain.Ticket, error) {
	tickets, err := s.repo.List(ticketOpts)

	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (s *ticketService) UpdateTicket(ticketId uuid.UUID, opts ports.TicketUpdateOptions) (*domain.Ticket, error) {
	ticket, err := s.repo.Update(ticketId, opts)

	if err != nil {
		return nil, err
	}

	return ticket, nil
}
