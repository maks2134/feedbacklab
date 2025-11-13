package user_projects

import "context"

type Service interface {
	Create(ctx context.Context, up *UserProject) error
	Get(ctx context.Context, userID string, projectID int) (*UserProject, error)
	GetAll(ctx context.Context) ([]UserProject, error)
	Update(ctx context.Context, up *UserProject) error
	Delete(ctx context.Context, userID string, projectID int) error
}

type userProjectService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &userProjectService{repo: repo}
}

func (s *userProjectService) Create(ctx context.Context, up *UserProject) error {
	return s.repo.Create(ctx, up)
}

func (s *userProjectService) Get(ctx context.Context, userID string, projectID int) (*UserProject, error) {
	return s.repo.Get(ctx, userID, projectID)
}

func (s *userProjectService) GetAll(ctx context.Context) ([]UserProject, error) {
	return s.repo.GetAll(ctx)
}

func (s *userProjectService) Update(ctx context.Context, up *UserProject) error {
	return s.repo.Update(ctx, up)
}

func (s *userProjectService) Delete(ctx context.Context, userID string, projectID int) error {
	return s.repo.Delete(ctx, userID, projectID)
}
