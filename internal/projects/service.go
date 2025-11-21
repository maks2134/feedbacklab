package projects

import (
	"context"
	"innotech/pkg/logger" // Импортируем ваш логгер
)

type Service interface {
	Create(ctx context.Context, p *Project) error
	GetByID(ctx context.Context, id int) (*Project, error)
	GetAll(ctx context.Context) ([]Project, error)
	Update(ctx context.Context, p *Project) error
	Delete(ctx context.Context, id int) error
}

type projectService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	logger.Info("project service initialized")
	return &projectService{repo: repo}
}

func (s *projectService) Create(ctx context.Context, p *Project) error {
	logger.Info("service: create project",
		"name", p.Name,
		"gitlab_project_id", p.GitlabProjectID,
	)

	if err := s.repo.Create(ctx, p); err != nil {
		logger.Error("service: create failed",
			"error", err.Error(),
			"name", p.Name,
		)
		return err
	}

	logger.Info("service: project created successfully",
		"id", p.ID,
		"name", p.Name,
	)
	return nil
}

func (s *projectService) GetByID(ctx context.Context, id int) (*Project, error) {
	logger.Debug("service: get by id", "id", id)

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("service: get by id failed",
			"id", id,
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Debug("service: project retrieved successfully",
		"id", p.ID,
		"name", p.Name,
	)
	return p, nil
}

func (s *projectService) GetAll(ctx context.Context) ([]Project, error) {
	logger.Debug("service: get all projects")

	ps, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.Error("service: get all failed",
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Info("service: projects list retrieved",
		"count", len(ps),
	)
	return ps, nil
}

func (s *projectService) Update(ctx context.Context, p *Project) error {
	logger.Info("service: update project",
		"id", p.ID,
		"name", p.Name,
	)

	if err := s.repo.Update(ctx, p); err != nil {
		logger.Error("service: update failed",
			"id", p.ID,
			"error", err.Error(),
		)
		return err
	}

	logger.Info("service: project updated successfully",
		"id", p.ID,
		"name", p.Name,
	)
	return nil
}

func (s *projectService) Delete(ctx context.Context, id int) error {
	logger.Warn("service: delete project", "id", id)

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("service: delete failed",
			"id", id,
			"error", err.Error(),
		)
		return err
	}

	logger.Info("service: project deleted successfully", "id", id)
	return nil
}
