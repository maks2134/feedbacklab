package service_test

import (
	"errors"
	"innotech/internal/models"
	"innotech/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDocumentationRepository struct{ mock.Mock }

func (m *MockDocumentationRepository) Create(p *models.Documentation) error { return m.Called(p).Error(0) }
func (m *MockDocumentationRepository) GetAll() ([]models.Documentation, error) {
	args := m.Called()
	return args.Get(0).([]models.Documentation), args.Error(1)
}
func (m *MockDocumentationRepository) GetByID(id int) (*models.Documentation, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Documentation), args.Error(1)
}
func (m *MockDocumentationRepository) Update(p *models.Documentation) error { return m.Called(p).Error(0) }
func (m *MockDocumentationRepository) Delete(id int) error                  { return m.Called(id).Error(0) }

func TestDocumentationService_CRUD(t *testing.T) {
	mockRepo := new(MockDocumentationRepository)
	svc := service.NewDocumentationService(mockRepo)

	doc := &models.Documentation{ID: 1, Title: "API Spec"}

	mockRepo.On("Create", doc).Return(nil)
	mockRepo.On("GetAll").Return([]models.Documentation{*doc}, nil)
	mockRepo.On("GetByID", 1).Return(doc, nil)
	mockRepo.On("Update", doc).Return(nil)
	mockRepo.On("Delete", 1).Return(nil)

	assert.NoError(t, svc.Create(doc))
	list, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	got, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "API Spec", got.Title)

	assert.NoError(t, svc.Update(doc))
	assert.NoError(t, svc.Delete(1))
}
