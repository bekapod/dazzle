package application

import (
	"context"

	"dazzle/internal/domain"
)

// SpecService implements domain.SpecService.
type SpecService struct {
	repo domain.SpecRepository
}

func NewSpecService(repo domain.SpecRepository) *SpecService {
	return &SpecService{repo: repo}
}

func (s *SpecService) LoadSpec(ctx context.Context, source string) (*domain.Spec, error) {
	return s.repo.Load(ctx, source)
}

func (s *SpecService) GetInfo(spec *domain.Spec) domain.SpecInfo {
	return spec.Info
}
