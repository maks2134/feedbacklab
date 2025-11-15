package service_test

import (
	"errors"
	"innotech/internal/models"
	"innotech/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockContractRepository struct{ mock.Mock }

func (m *MockContractRepository) Create(p *models.Contract) error { return m.Called(p).Error(0) }
func (m *MockContractRepository) GetAll() ([]models.Contract, error) {
	args := m.Called()
	return args.Get(0).([]models.Contract), args.Error(1)
}
func (m *MockContractRepository) GetByID(id int) (*models.Contract, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Contract), args.Error(1)
}
func (m *MockContractRepository) Update(p *models.Contract) error { return m.Called(p).Error(0) }
func (m *MockContractRepository) Delete(id int) error             { return m.Called(id).Error(0) }

func TestContractService_CRUD(t *testing.T) {
	mockRepo := new(MockContractRepository)
	svc := service.NewContractService(mockRepo)

	contract := &models.Contract{ID: 1, Title: "Agreement"}

	mockRepo.On("Create", contract).Return(nil)
	mockRepo.On("GetAll").Return([]models.Contract{*contract}, nil)
	mockRepo.On("GetByID", 1).Return(contract, nil)
	mockRepo.On("Update", contract).Return(nil)
	mockRepo.On("Delete", 1).Return(nil)

	assert.NoError(t, svc.Create(contract))
	items, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, items, 1)

	got, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Agreement", got.Title)

	assert.NoError(t, svc.Update(contract))
	assert.NoError(t, svc.Delete(1))
}
