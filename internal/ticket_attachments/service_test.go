package ticket_attachments

import (
	"errors"
	"innotech/internal/storage/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepoTA struct{ mock.Mock }

func (m *mockRepoTA) Create(att *postgres.TicketAttachment) error { return m.Called(att).Error(0) }
func (m *mockRepoTA) GetByID(id int) (*postgres.TicketAttachment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postgres.TicketAttachment), args.Error(1)
}
func (m *mockRepoTA) GetByTicketID(ticketID int) ([]postgres.TicketAttachment, error) {
	args := m.Called(ticketID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgres.TicketAttachment), args.Error(1)
}
func (m *mockRepoTA) Update(att *postgres.TicketAttachment) error { return m.Called(att).Error(0) }
func (m *mockRepoTA) Delete(id int) error                         { return m.Called(id).Error(0) }

func TestCreate_ValidationFails_WhenEmptyFilePath(t *testing.T) {
	repo := new(mockRepoTA)
	svc := NewService(repo)

	err := svc.Create(&postgres.TicketAttachment{FilePath: ""})
	assert.Error(t, err)
}

func TestCreate_HappyPath_CallsRepo(t *testing.T) {
	repo := new(mockRepoTA)
	svc := NewService(repo)

	att := &postgres.TicketAttachment{TicketID: 1, FilePath: "p"}
	repo.On("Create", att).Return(nil).Once()

	err := svc.Create(att)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreate_RepoError_ReturnsError(t *testing.T) {
	repo := new(mockRepoTA)
	svc := NewService(repo)

	att := &postgres.TicketAttachment{TicketID: 1, FilePath: "p"}
	repo.On("Create", att).Return(errors.New("db")).Once()

	err := svc.Create(att)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestGetAndUpdateDelete_PassesThrough(t *testing.T) {
	repo := new(mockRepoTA)
	svc := NewService(repo)

	exp := &postgres.TicketAttachment{ID: 10}
	repo.On("GetByID", 10).Return(exp, nil).Once()
	repo.On("GetByTicketID", 2).Return([]postgres.TicketAttachment{{ID: 1}}, nil).Once()
	repo.On("Update", exp).Return(nil).Once()
	repo.On("Delete", 10).Return(nil).Once()

	got, err := svc.GetByID(10)
	assert.NoError(t, err)
	assert.Equal(t, exp, got)

	list, err := svc.GetByTicketID(2)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	assert.NoError(t, svc.Update(exp))
	assert.NoError(t, svc.Delete(10))
	repo.AssertExpectations(t)
}
