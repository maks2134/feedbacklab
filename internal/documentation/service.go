package documentation

type DocumentationService struct {
	repo *DocumentationRepository
}

func NewDocumentationService(repo *DocumentationRepository) *DocumentationService {
	return &DocumentationService{repo: repo}
}

func (s *DocumentationService) GetAll() ([]Documentation, error) {
	return s.repo.GetAll()
}

func (s *DocumentationService) GetByID(id int) (*Documentation, error) {
	return s.repo.GetByID(id)
}

func (s *DocumentationService) Create(d *Documentation) error {
	return s.repo.Create(d)
}

func (s *DocumentationService) Update(d *Documentation) error {
	return s.repo.Update(d)
}

func (s *DocumentationService) Delete(id int) error {
	return s.repo.Delete(id)
}
