package domain

import (
	"time"

	"github.com/google/uuid"
)

type TicketStatus string

const (
	TicketStatusNew        TicketStatus = "NEW"
	TicketStatusInProgress TicketStatus = "IN_PROGRESS"
	TicketStatusResolved   TicketStatus = "RESOLVED"
)

type Ticket struct {
	ID          int64        `json:"id"`
	TicketID    uuid.UUID    `json:"ticket_id"`
	Name        string       `json:"name,omitempty"`
	IssuerEmail string       `json:"issuer_email,omitempty"`
	Description *string      `json:"description,omitempty"`
	Status      TicketStatus `json:"status,omitempty"`
	Response    *string      `json:"response,omitempty"`
	Note        *string      `json:"note,omitempty"`
	CreatedAt   *time.Time   `json:"created_at,omitempty"`
	UpdatedAt   *time.Time   `json:"updated_at,omitempty"`
	// Code        string       `json:"code"`
}

func (ts TicketStatus) IsValid() bool {
	switch ts {
	case TicketStatusNew, TicketStatusInProgress, TicketStatusResolved:
		return true
	}
	return false
}
