package service_test

import (
	"errors"
	"innotech/internal/models"
	"innotech/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockModuleRepository struct{ mock.Mock }

func (m *MockModuleRepository) Create(p *models.Module) error { return m.Called(p).Error(0) }
func (m *MockModuleRepository) GetAll() ([]models.Module, error) {
	args := m.Called()
	return args.Get(0).([]models.Module), args.Error(1)
}
func (m *MockModuleRepository) GetByID(id int) (*models.Module, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Module), args.Error(1)
}
func (m *MockModuleRepository) Update(p *models.Module) error { return m.Called(p).Error(0) }
func (m *MockModuleRepository) Delete(id int) error           { return m.Called(id).Error(0) }

func TestModuleService_CRUD(t *testing.T) {
	mockRepo := new(MockModuleRepository)
	svc := service.NewModuleService(mockRepo)

	module := &models.Module{ID: 1, Name: "Core"}

	mockRepo.On("Create", module).Return(nil)
	mockRepo.On("GetAll").Return([]models.Module{*module}, nil)
	mockRepo.On("GetByID", 1).Return(module, nil)
	mockRepo.On("Update", module).Return(nil)
	mockRepo.On("Delete", 1).Return(nil)

	assert.NoError(t, svc.Create(module))
	list, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	got, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Core", got.Name)

	assert.NoError(t, svc.Update(module))
	assert.NoError(t, svc.Delete(1))
}
