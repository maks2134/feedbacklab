package ticket_chats

import (
	"errors"
	"innotech/internal/storage/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Create(chat *postgres.TicketChat) error {
	args := m.Called(chat)
	return args.Error(0)
}
func (m *mockRepo) GetByID(id int) (*postgres.TicketChat, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postgres.TicketChat), args.Error(1)
}
func (m *mockRepo) GetByTicketID(ticketID int) ([]postgres.TicketChat, error) {
	args := m.Called(ticketID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgres.TicketChat), args.Error(1)
}
func (m *mockRepo) Update(chat *postgres.TicketChat) error {
	args := m.Called(chat)
	return args.Error(0)
}
func (m *mockRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreate_ValidationFails_WhenEmptyMessage(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	err := svc.Create(&postgres.TicketChat{Message: ""})
	assert.Error(t, err)
}

func TestCreate_HappyPath_CallsRepo(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	chat := &postgres.TicketChat{TicketID: 1, SenderID: "u", Message: "hi"}
	repo.On("Create", chat).Return(nil).Once()

	err := svc.Create(chat)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreate_RepoError_ReturnsError(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	chat := &postgres.TicketChat{TicketID: 1, SenderID: "u", Message: "hi"}
	repo.On("Create", chat).Return(errors.New("db")).Once()

	err := svc.Create(chat)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestGetByID_PassesThrough(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	exp := &postgres.TicketChat{ID: 2}
	repo.On("GetByID", 2).Return(exp, nil).Once()

	got, err := svc.GetByID(2)
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
	repo.AssertExpectations(t)
}

func TestGetByTicketID_PassesThrough(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	list := []postgres.TicketChat{{ID: 1}, {ID: 2}}
	repo.On("GetByTicketID", 5).Return(list, nil).Once()

	got, err := svc.GetByTicketID(5)
	assert.NoError(t, err)
	assert.Equal(t, list, got)
	repo.AssertExpectations(t)
}

func TestUpdate_Delete_PassThrough(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)

	c := &postgres.TicketChat{ID: 3, Message: "ok"}
	repo.On("Update", c).Return(nil).Once()
	repo.On("Delete", 3).Return(nil).Once()

	assert.NoError(t, svc.Update(c))
	assert.NoError(t, svc.Delete(3))
	repo.AssertExpectations(t)
}
