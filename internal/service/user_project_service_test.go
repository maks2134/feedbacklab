package service_test

import (
	"errors"
	"innotech/internal/models"
	"innotech/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserProjectRepository struct{ mock.Mock }

func (m *MockUserProjectRepository) Create(p *models.UserProject) error { return m.Called(p).Error(0) }
func (m *MockUserProjectRepository) GetAll() ([]models.UserProject, error) {
	args := m.Called()
	return args.Get(0).([]models.UserProject), args.Error(1)
}
func (m *MockUserProjectRepository) GetByID(id int) (*models.UserProject, error) {
	args := m.Called(id)
	return args.Get(0).(*models.UserProject), args.Error(1)
}
func (m *MockUserProjectRepository) Update(p *models.UserProject) error { return m.Called(p).Error(0) }
func (m *MockUserProjectRepository) Delete(id int) error                { return m.Called(id).Error(0) }

func TestUserProjectService_CRUD(t *testing.T) {
	mockRepo := new(MockUserProjectRepository)
	svc := service.NewUserProjectService(mockRepo)

	relation := &models.UserProject{ID: 1, UserID: "uuid-1", ProjectID: 2}

	mockRepo.On("Create", relation).Return(nil)
	mockRepo.On("GetAll").Return([]models.UserProject{*relation}, nil)
	mockRepo.On("GetByID", 1).Return(relation, nil)
	mockRepo.On("Update", relation).Return(nil)
	mockRepo.On("Delete", 1).Return(nil)

	assert.NoError(t, svc.Create(relation))
	list, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	got, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "uuid-1", got.UserID)

	assert.NoError(t, svc.Update(relation))
	assert.NoError(t, svc.Delete(1))
}
