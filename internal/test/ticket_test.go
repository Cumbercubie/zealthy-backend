package test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zealthy/helpdesk-backend/internal/core/domain"
	"github.com/zealthy/helpdesk-backend/internal/core/ports"
	"github.com/zealthy/helpdesk-backend/internal/core/services"
)

type MockTicketRepo struct {
	mock.Mock
}

type MockTicketDBHandler struct {
	mock.Mock
}

func (r *MockTicketRepo) Create(ticketDetail ports.TicketCreateDetail) (*domain.Ticket, error) {
	args := r.Called(ticketDetail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (r *MockTicketRepo) GetByID(id uuid.UUID) (*domain.Ticket, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (r *MockTicketRepo) Update(ticketId uuid.UUID, options ports.TicketUpdateOptions) (*domain.Ticket, error) {
	args := r.Called(ticketId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (r *MockTicketRepo) List(ticketOpts *ports.TicketListOptions) ([]*domain.Ticket, error) {
	args := r.Called(ticketOpts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Ticket), args.Error(1)
}

func TestNewTicketService(t *testing.T) {
	mockTicketRepo := &MockTicketRepo{}

	ticketService := services.NewTicketService(mockTicketRepo)

	assert.NotNil(t, ticketService)
	assert.Equal(t, mockTicketRepo, ticketService.Repo)
}

func TestGetListTicketSuccess(t *testing.T) {
	var responseTicket *domain.Ticket
	mockTicketRepo := &MockTicketRepo{}

	ticketService := services.NewTicketService(mockTicketRepo)
	testTime := time.Now()
	testId := uuid.New()
	testEmail := "test@gmail.com"
	testString := "testString"
	testStatus := domain.TicketStatusNew

	responseTicket = &domain.Ticket{
		ID:          0,
		TicketID:    testId,
		Name:        testString,
		IssuerEmail: testEmail,
		Description: &testString,
		Status:      testStatus,
		Response:    &testString,
		Note:        &testString,
		CreatedAt:   &testTime,
		UpdatedAt:   &testTime,
	}

	mockTicketRepo.On("List", mock.Anything).Return([]*domain.Ticket{responseTicket}, nil)

	response, err := ticketService.ListTickets(nil)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response[0].TicketID, testId)
	assert.Equal(t, response[0].Name, testString)
	assert.Equal(t, *response[0].Description, testString)
	assert.Equal(t, response[0].IssuerEmail, testEmail)
	assert.Equal(t, response[0].Status, testStatus)
	assert.Equal(t, *response[0].CreatedAt, testTime)
	assert.Equal(t, *response[0].UpdatedAt, testTime)
}

func TestGetTicketById(t *testing.T) {
	var responseTicket *domain.Ticket
	mockTicketRepo := &MockTicketRepo{}

	ticketService := services.NewTicketService(mockTicketRepo)
	testTime := time.Now()
	testId := uuid.New()
	testEmail := "test@gmail.com"
	testString := "testString"
	testStatus := domain.TicketStatusNew

	responseTicket = &domain.Ticket{
		ID:          0,
		TicketID:    testId,
		Name:        testString,
		IssuerEmail: testEmail,
		Description: &testString,
		Status:      testStatus,
		Response:    &testString,
		Note:        &testString,
		CreatedAt:   &testTime,
		UpdatedAt:   &testTime,
	}

	mockTicketRepo.On("GetByID", mock.Anything).Return(responseTicket, nil)

	response, err := ticketService.GetTicket(testId)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.TicketID, testId)
	assert.Equal(t, response.Name, testString)
	assert.Equal(t, *response.Description, testString)
	assert.Equal(t, response.IssuerEmail, testEmail)
	assert.Equal(t, response.Status, testStatus)
	assert.Equal(t, *response.CreatedAt, testTime)
	assert.Equal(t, *response.UpdatedAt, testTime)
}

func TestGetTicketByIdWrongId(t *testing.T) {
	// var responseTicket *domain.Ticket
	mockTicketRepo := &MockTicketRepo{}

	ticketService := services.NewTicketService(mockTicketRepo)
	// testTime := time.Now()
	testId := uuid.New()
	// testEmail := "test@gmail.com"

	mockTicketRepo.On("GetByID", mock.Anything).Return(nil, errors.New("not found ticket with this id"))

	response, err := ticketService.GetTicket(testId)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "not found ticket with this id")
	assert.Nil(t, response)
}

func TestCreateTicketSuccess(t *testing.T) {
	var responseTicket *domain.Ticket
	mockTicketRepo := &MockTicketRepo{}

	ticketService := services.NewTicketService(mockTicketRepo)
	testTime := time.Now()
	testId := uuid.New()
	testEmail := "test@gmail.com"

	testString := "testString"
	testStatus := domain.TicketStatusNew

	responseTicket = &domain.Ticket{
		ID:          0,
		TicketID:    testId,
		Name:        testString,
		IssuerEmail: testEmail,
		Description: &testString,
		Status:      testStatus,
		Response:    &testString,
		Note:        &testString,
		CreatedAt:   &testTime,
		UpdatedAt:   &testTime,
	}
	mockTicketRepo.On("Create", mock.Anything).Return(responseTicket, nil)

	response, err := ticketService.CreateTicket(ports.TicketCreateDetail{
		Name:        testString,
		IssuerEmail: testEmail,
		Description: &testString,
	})

	assert.Nil(t, err)
	assert.NotNil(t, response)

}

func TestUpdateTicketSuccess(t *testing.T) {
	var responseTicket *domain.Ticket
	mockTicketRepo := &MockTicketRepo{}

	ticketService := services.NewTicketService(mockTicketRepo)
	testTime := time.Now()
	testId := uuid.New()
	testEmail := "test@gmail.com"

	testString := "testString"
	testStatus := domain.TicketStatusNew

	responseTicket = &domain.Ticket{
		ID:          0,
		TicketID:    testId,
		Name:        testString,
		IssuerEmail: testEmail,
		Description: &testString,
		Status:      testStatus,
		Response:    &testString,
		Note:        &testString,
		CreatedAt:   &testTime,
		UpdatedAt:   &testTime,
	}
	mockTicketRepo.On("Update", mock.Anything).Return(responseTicket, nil)

	response, err := ticketService.UpdateTicket(testId, ports.TicketUpdateOptions{
		Status:   &testStatus,
		Response: &testString,
		Note:     &testString,
	})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.TicketID, testId)
	assert.Equal(t, response.Name, testString)
	assert.Equal(t, *response.Description, testString)
	assert.Equal(t, response.IssuerEmail, testEmail)
	assert.Equal(t, response.Status, testStatus)
	assert.Equal(t, *response.CreatedAt, testTime)
	assert.Equal(t, *response.UpdatedAt, testTime)

}

func TestUpdateTicketError(t *testing.T) {
	mockTicketRepo := &MockTicketRepo{}

	ticketService := services.NewTicketService(mockTicketRepo)
	testId := uuid.New()

	testString := "testString"
	testStatus := domain.TicketStatusInProgress

	mockTicketRepo.On("Update", mock.Anything).Return(nil, errors.New("test error"))

	response, err := ticketService.UpdateTicket(testId, ports.TicketUpdateOptions{
		Status:   &testStatus,
		Response: &testString,
		Note:     &testString,
	})

	assert.NotNil(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "test error")
}
