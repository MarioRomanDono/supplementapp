package supplement

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("supplement not found")
	ErrAlreadyExists = errors.New("supplement already exists")
)

type SupplementService struct {
	repository SupplementRepository
}

func NewSupplementService(repository SupplementRepository) *SupplementService {
	return &SupplementService{repository: repository}
}

func (service *SupplementService) Create(supplement Supplement) error {
	existing, err := service.repository.FindByGtin(supplement.Gtin)

	if err != nil {
		return err
	}

	if existing != nil {
		return fmt.Errorf("%v: %w", supplement, ErrAlreadyExists)
	}

	return service.repository.Create(supplement)
}

func (service *SupplementService) FindByGtin(gtin string) (*Supplement, error) {
	supplement, err := service.repository.FindByGtin(gtin)

	if err != nil {
		return nil, err
	}

	if supplement == nil {
		return nil, fmt.Errorf("%s: %w", gtin, ErrNotFound)
	}

	return supplement, nil
}

func (service *SupplementService) Delete(gtin string) error {
	supplement, err := service.repository.FindByGtin(gtin)

	if err != nil {
		return err
	}

	if supplement == nil {
		return fmt.Errorf("%s: %w", gtin, ErrNotFound)
	}

	return service.repository.Delete(*supplement)
}
