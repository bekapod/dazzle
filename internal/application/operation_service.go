package application

import (
	"sort"
	"strings"

	"dazzle/internal/domain"
)

var methodOrder = map[domain.HTTPMethod]int{
	domain.GET:     1,
	domain.POST:    2,
	domain.PUT:     3,
	domain.PATCH:   4,
	domain.DELETE:  5,
	domain.HEAD:    6,
	domain.OPTIONS: 7,
}

// OperationService implements domain.OperationService.
type OperationService struct{}

func NewOperationService() *OperationService {
	return &OperationService{}
}

func (s *OperationService) ListOperations(spec *domain.Spec) []domain.Operation {
	return spec.Operations
}

func (s *OperationService) FilterOperations(operations []domain.Operation, filter domain.OperationFilter) []domain.Operation {
	var result []domain.Operation
	for _, op := range operations {
		if matchesFilter(op, filter) {
			result = append(result, op)
		}
	}
	return result
}

func (s *OperationService) SortOperations(operations []domain.Operation) []domain.Operation {
	sorted := make([]domain.Operation, len(operations))
	copy(sorted, operations)

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Path != sorted[j].Path {
			return sorted[i].Path < sorted[j].Path
		}
		return methodOrder[sorted[i].Method] < methodOrder[sorted[j].Method]
	})

	return sorted
}

func matchesFilter(op domain.Operation, f domain.OperationFilter) bool {
	if f.Query != "" {
		q := strings.ToLower(f.Query)
		if !strings.Contains(strings.ToLower(op.Path), q) &&
			!strings.Contains(strings.ToLower(op.Summary), q) {
			return false
		}
	}

	if f.Method != "" && op.Method != f.Method {
		return false
	}

	if len(f.Tags) > 0 && !hasMatchingTag(op.Tags, f.Tags) {
		return false
	}

	return true
}

func hasMatchingTag(opTags, filterTags []string) bool {
	for _, ft := range filterTags {
		for _, ot := range opTags {
			if ft == ot {
				return true
			}
		}
	}
	return false
}
