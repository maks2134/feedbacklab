package service_test

import (
	"errors"
	"innotech/internal/models"
	"innotech/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProjectRepository struct{ mock.Mock }

func (m *MockProjectRepository) Create(p *models.Project) error { return m.Called(p).Error(0) }
func (m *MockProjectRepository) GetAll() ([]models.Project, error) {
	args := m.Called()
	return args.Get(0).([]models.Project), args.Error(1)
}
func (m *MockProjectRepository) GetByID(id int) (*models.Project, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Project), args.Error(1)
}
func (m *MockProjectRepository) Update(p *models.Project) error { return m.Called(p).Error(0) }
func (m *MockProjectRepository) Delete(id int) error            { return m.Called(id).Error(0) }

func TestProjectService_CRUD(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	svc := service.NewProjectService(mockRepo)

	project := &models.Project{ID: 1, Name: "Demo"}

	mockRepo.On("Create", project).Return(nil)
	mockRepo.On("GetAll").Return([]models.Project{*project}, nil)
	mockRepo.On("GetByID", 1).Return(project, nil)
	mockRepo.On("Update", project).Return(nil)
	mockRepo.On("Delete", 1).Return(nil)

	assert.NoError(t, svc.Create(project))
	items, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, items, 1)

	got, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Demo", got.Name)

	assert.NoError(t, svc.Update(project))
	assert.NoError(t, svc.Delete(1))
}
