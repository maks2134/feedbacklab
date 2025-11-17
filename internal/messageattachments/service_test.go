package messageattachments

import (
	"errors"
	"innotech/internal/storage/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepoMA struct{ mock.Mock }

func (m *mockRepoMA) Create(att *postgres.MessageAttachment) error { return m.Called(att).Error(0) }
func (m *mockRepoMA) GetByID(id int) (*postgres.MessageAttachment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postgres.MessageAttachment), args.Error(1)
}
func (m *mockRepoMA) GetByChatID(chatID int) ([]postgres.MessageAttachment, error) {
	args := m.Called(chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgres.MessageAttachment), args.Error(1)
}
func (m *mockRepoMA) Update(att *postgres.MessageAttachment) error { return m.Called(att).Error(0) }
func (m *mockRepoMA) Delete(id int) error                          { return m.Called(id).Error(0) }

func TestCreate_ValidationFails_WhenMissingFields(t *testing.T) {
	repo := new(mockRepoMA)
	svc := NewService(repo)

	// missing file path
	err := svc.Create(&postgres.MessageAttachment{ChatID: 1, FilePath: ""})
	assert.Error(t, err)

	// missing chat id
	err = svc.Create(&postgres.MessageAttachment{ChatID: 0, FilePath: "p"})
	assert.Error(t, err)
}

func TestCreate_HappyPath_CallsRepo(t *testing.T) {
	repo := new(mockRepoMA)
	svc := NewService(repo)

	a := &postgres.MessageAttachment{ChatID: 2, FilePath: "p"}
	repo.On("Create", a).Return(nil).Once()

	err := svc.Create(a)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreate_RepoError_ReturnsError(t *testing.T) {
	repo := new(mockRepoMA)
	svc := NewService(repo)

	a := &postgres.MessageAttachment{ChatID: 2, FilePath: "p"}
	repo.On("Create", a).Return(errors.New("db")).Once()

	err := svc.Create(a)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestGetUpdateDelete_PassesThrough(t *testing.T) {
	repo := new(mockRepoMA)
	svc := NewService(repo)

	exp := &postgres.MessageAttachment{ID: 11}
	repo.On("GetByID", 11).Return(exp, nil).Once()
	repo.On("GetByChatID", 3).Return([]postgres.MessageAttachment{{ID: 1}}, nil).Once()
	repo.On("Update", exp).Return(nil).Once()
	repo.On("Delete", 11).Return(nil).Once()

	got, err := svc.GetByID(11)
	assert.NoError(t, err)
	assert.Equal(t, exp, got)

	list, err := svc.GetByChatID(3)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	assert.NoError(t, svc.Update(exp))
	assert.NoError(t, svc.Delete(11))
	repo.AssertExpectations(t)
}
