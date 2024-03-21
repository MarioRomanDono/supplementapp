package supplement

import (
	"context"
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

func (service *SupplementService) Create(ctx context.Context, supplement Supplement) error {
	existing, err := service.repository.FindByGtin(ctx, supplement.Gtin)

	if err != nil {
		return err
	}

	if existing != nil {
		return fmt.Errorf("%v: %w", supplement, ErrAlreadyExists)
	}

	return service.repository.Create(ctx, supplement)
}

func (service *SupplementService) FindByGtin(ctx context.Context, gtin string) (*Supplement, error) {
	supplement, err := service.repository.FindByGtin(ctx, gtin)

	if err != nil {
		return nil, err
	}

	if supplement == nil {
		return nil, fmt.Errorf("%s: %w", gtin, ErrNotFound)
	}

	return supplement, nil
}

func (service *SupplementService) Delete(ctx context.Context, gtin string) error {
	supplement, err := service.repository.FindByGtin(ctx, gtin)

	if err != nil {
		return err
	}

	if supplement == nil {
		return fmt.Errorf("%s: %w", gtin, ErrNotFound)
	}

	return service.repository.Delete(ctx, *supplement)
}

func (service *SupplementService) Update(ctx context.Context, gtin string, other UpdatableSupplement) error {
	supplement, err := service.repository.FindByGtin(ctx, gtin)

	if err != nil {
		return err
	}

	if supplement == nil {
		return fmt.Errorf("%s: %w", gtin, ErrNotFound)
	}

	updated := supplement.update(other)

	return service.repository.Update(ctx, updated)
}

// TODO: Add pagination, sorting, and filtering
func (service *SupplementService) ListAll(ctx context.Context) ([]Supplement, error) {
	return service.repository.ListAll(ctx)
}
