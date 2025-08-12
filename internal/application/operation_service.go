package application

import (
	"sort"

	"dazzle/internal/domain"
)

type OperationServiceImpl struct {
	repo domain.OperationRepository
}

func NewOperationService(repo domain.OperationRepository) domain.OperationService {
	return &OperationServiceImpl{
		repo: repo,
	}
}

func (s *OperationServiceImpl) ListOperations() ([]domain.Operation, error) {
	operations, err := s.repo.GetOperations()
	if err != nil {
		return nil, err
	}
	return s.SortOperations(operations), nil
}





func (s *OperationServiceImpl) SortOperations(operations []domain.Operation) []domain.Operation {
	methodOrder := map[domain.HTTPMethod]int{
		domain.GET:    1,
		domain.POST:   2,
		domain.PUT:    3,
		domain.PATCH:  4,
		domain.DELETE: 5,
	}

	sort.Slice(operations, func(i, j int) bool {
		if operations[i].Path == operations[j].Path {
			return methodOrder[operations[i].Method] < methodOrder[operations[j].Method]
		}
		return operations[i].Path < operations[j].Path
	})

	return operations
}