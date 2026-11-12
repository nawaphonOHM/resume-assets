package handler

import (
	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

func requestAssetCriteria(c rest.Context) (*string, error) {
	if c == nil {
		return nil, nil
	}
	query := make(map[string][]string)
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, err
	}
	return criteriaFromQuery(query)
}

func criteriaFromQuery(query map[string][]string) (*string, error) {
	if len(query) == 0 {
		return nil, nil
	}
	if len(query) != 1 {
		return nil, domain.ErrInvalidAssetCriteria
	}
	return singleCriteriaQuery(query)
}

func singleCriteriaQuery(query map[string][]string) (*string, error) {
	for criteria, values := range query {
		return validateSingleCriteria(criteria, values)
	}
	return nil, domain.ErrInvalidAssetCriteria
}

func validateSingleCriteria(criteria string, values []string) (*string, error) {
	if len(values) != 1 {
		return nil, domain.ErrInvalidAssetCriteria
	}
	if values[0] != "" {
		return nil, domain.ErrInvalidAssetCriteria
	}
	return &criteria, nil
}
