package health

import "context"

type SelfHealthService struct{}

func NewSelfHealthService() *SelfHealthService {
	return &SelfHealthService{}
}

func (s *SelfHealthService) Check(ctx context.Context) error {
	return nil
}
