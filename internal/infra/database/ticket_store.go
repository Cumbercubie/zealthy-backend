package sql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/zealthy/helpdesk-backend/internal/core/domain"
	"github.com/zealthy/helpdesk-backend/internal/core/ports"
)

const (
	CREATE_TICKET_QUERY          = "INSERT INTO tickets (ticket_id, name, issuer_email, description, status, response, note, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, ticket_id, name, issuer_email, description, status, response, note, created_at, updated_at"
	LIST_TICKETS_QUERY           = "SELECT id, ticket_id, name, issuer_email, description, status, response, note, created_at, updated_at FROM tickets"
	LIST_TICKETS_QUERY_BY_STATUS = "SELECT id, ticket_id, name, issuer_email, description, status, response, note, created_at, updated_at FROM tickets WHERE status = $1"
	GET_TICKET_BY_ID_QUERY       = "SELECT id, ticket_id, name, issuer_email, description, status, response, note, created_at, updated_at FROM tickets WHERE ticket_id = $1"
	UPDATE_TICKET_QUERY          = "UPDATE tickets SET status = $1, response = $2, note = $3 WHERE ticket_id = $4 RETURNING id, ticket_id, name, issuer_email, description, status, response, note, created_at, updated_at"
)

func (rh *DBHandler) CreateTicket(detail ports.TicketCreateDetail) (*domain.Ticket, error) {
	dbConn, err := rh.getDbConnection()

	if err != nil {
		return nil, err
	}

	defer dbConn.Release()

	tx, err := dbConn.Begin(context.Background())

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	var ticket *domain.Ticket = &domain.Ticket{}
	currentDate := time.Now()
	genTicketInfo(ticket)
	err = tx.QueryRow(context.Background(), CREATE_TICKET_QUERY,
		&ticket.TicketID,
		&detail.Name,
		&detail.IssuerEmail,
		&detail.Description,
		domain.TicketStatusNew,
		&detail.Response,
		&detail.Note,
		&currentDate,
		&currentDate,
	).Scan(
		&ticket.ID,
		&ticket.TicketID,
		&ticket.Name,
		&ticket.IssuerEmail,
		&ticket.Description,
		&ticket.Status,
		&ticket.Response,
		&ticket.Note,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket: ", err.Error())
	}

	return ticket, nil
}

func (rh *DBHandler) GetTicketById(id uuid.UUID) (*domain.Ticket, error) {
	dbConn, err := rh.getDbConnection()
	if err != nil {
		return nil, err
	}
	defer dbConn.Release()
	ticket := &domain.Ticket{}
	err = dbConn.QueryRow(context.Background(), GET_TICKET_BY_ID_QUERY, id).Scan(
		&ticket.ID,
		&ticket.TicketID,
		&ticket.Name,
		&ticket.IssuerEmail,
		&ticket.Description,
		&ticket.Status,
		&ticket.Response,
		&ticket.Note,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (h *DBHandler) ListTickets(ticketOpts *ports.TicketListOptions) ([]*domain.Ticket, error) {
	dbConn, err := h.getDbConnection()
	var rows pgx.Rows
	if err != nil {
		return nil, err
	}
	defer dbConn.Release()
	query, args := buildListTicketQuery(ticketOpts)
	if len(query) <= 0 || args == nil {
		rows, err = dbConn.Query(
			context.Background(),
			LIST_TICKETS_QUERY,
		)
	} else {
		rows, err = dbConn.Query(
			context.Background(),
			query,
			args...,
		)
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tickets []*domain.Ticket

	for rows.Next() {

		var ticket domain.Ticket

		err := rows.Scan(
			&ticket.ID,
			&ticket.TicketID,
			&ticket.Name,
			&ticket.IssuerEmail,
			&ticket.Description,
			&ticket.Status,
			&ticket.Response,
			&ticket.Note,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		tickets = append(tickets, &ticket)

	}
	log.Println("tickets", len(tickets))

	if err := rows.Err(); err != nil {

		return nil, err
	}
	return tickets, nil
}

func (h *DBHandler) UpdateTicket(ticketId uuid.UUID, opts ports.TicketUpdateOptions) (*domain.Ticket, error) {
	dbConn, err := h.getDbConnection()
	if err != nil {
		return nil, err
	}
	defer dbConn.Release()

	updatedTicket := &domain.Ticket{}

	baseQuery, args := buildUpdateTicketQuery(&opts)

	if len(args) <= 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	err = dbConn.QueryRow(context.Background(), baseQuery,
		append(args, ticketId)...,
	).Scan(
		&updatedTicket.ID,
		&updatedTicket.TicketID,
		&updatedTicket.Name,
		&updatedTicket.IssuerEmail,
		&updatedTicket.Description,
		&updatedTicket.Status,
		&updatedTicket.Response,
		&updatedTicket.Note)

	if err != nil {
		return nil, err
	}

	return updatedTicket, nil

}

func genTicketInfo(ticket *domain.Ticket) {
	ticket.TicketID = uuid.New()
}

func buildListTicketQuery(options *ports.TicketListOptions) (string, []interface{}) {
	if options == nil {
		return "", nil
	}
	baseQuery := "SELECT * FROM tickets WHERE 1=1"
	var args []interface{}
	var conditions []string

	if options.Status != nil && *options.Status != "" {
		conditions = append(conditions, "status = $"+fmt.Sprint(len(args)+1))
		args = append(args, options.Status)
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + fmt.Sprint(conditions[0])
		for i := 1; i < len(conditions); i++ {
			baseQuery += " AND " + conditions[i]
		}
	}

	return baseQuery, args
}

func buildUpdateTicketQuery(options *ports.TicketUpdateOptions) (string, []interface{}) {
	if options == nil {
		return "", nil
	}
	baseQuery := "UPDATE tickets SET "
	var args []interface{}
	var conditions []string
	fmt.Println("options", options)
	if options.Status != nil && *options.Status != "" {
		conditions = append(conditions, "status = $"+fmt.Sprint(len(args)+1))
		args = append(args, options.Status)
	}

	if options.Response != nil && len(*options.Response) >= 0 {
		conditions = append(conditions, "response = $"+fmt.Sprint(len(args)+1))
		args = append(args, options.Response)
	}

	if options.Note != nil && len(*options.Note) >= 0 {
		conditions = append(conditions, "note = $"+fmt.Sprint(len(args)+1))
		args = append(args, options.Note)
	}

	if len(conditions) > 0 {
		baseQuery += fmt.Sprint(conditions[0])
		for i := 1; i < len(conditions); i++ {
			baseQuery += " AND " + conditions[i]
		}
	}

	baseQuery += fmt.Sprintf(" WHERE ticket_id = $%d RETURNING id, ticket_id, name, issuer_email, description, status, response, note", len(args)+1)

	return baseQuery, args
}
