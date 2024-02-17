package supplement

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrNotFound      = errors.New("supplement not found")
	ErrAlreadyExists = errors.New("supplement already exists")
	ErrInvalidFilterField  = errors.New("invalid filter field")
	ErrInvalidFilterOperator = errors.New("invalid filter operator")
	ErrInvalidFilterValue = errors.New("invalid filter value")
	ErrInvalidOrderField = errors.New("invalid order field")
	ErrInvalidOrderDirection = errors.New("invalid order direction")
	ErrInvalidLimit = errors.New("invalid limit")
	ErrInvalidCursorDirection = errors.New("invalid cursor direction")
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

func (service *SupplementService) Update(gtin string, other UpdatableSupplement) error {
	supplement, err := service.repository.FindByGtin(gtin)

	if err != nil {
		return err
	}

	if supplement == nil {
		return fmt.Errorf("%s: %w", gtin, ErrNotFound)
	}

	updated := supplement.update(other)

	return service.repository.Update(updated)
}

func (service *SupplementService) Search(options SearchOptions) (*SearchResponse, error) {
	err := validateSearchFilters(options.Filters)

	if err != nil {
		return nil, err
	}

	if options.Order != nil {
		err = validateSearchOrder(*options.Order)

		if err != nil {
			return nil, err
		}
	}

	err = validateLimit(options.Limit)

	if err != nil {
		return nil, err
	}

	if options.Cursor != nil {
		err = validateCursorDirection(*options.Cursor)

		if err != nil {
			return nil, err
		}
	}

	return service.repository.Search(options.Filters, options.Order, options.Limit, options.Cursor)
}

func validateSearchFilters(filters []SearchFilter) error {
	for _, filter := range filters {
		switch filter.Field {
		case "name", "brand", "flavor":
			err := validateStringFilter(filter)

			if err != nil {
				return err
			}
		case "carbohydrates", "electrolytes", "maltodextrose", "fructose", "caffeine", "sodium", "protein":
			err := validateFloatFilter(filter)

			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: %w", filter.Field, ErrInvalidFilterField)
		}
	}

	return nil
}

func validateStringFilter(filter SearchFilter) error {
	switch filter.Operator {
	case "=", "!=", "contains":
		return nil
	default:
		return fmt.Errorf("%s: %w", filter.Operator, ErrInvalidFilterOperator)
	}
}

func validateFloatFilter(filter SearchFilter) error {
	switch filter.Operator {
	case "=", "!=", ">", "<", ">=", "<=":
		_, err := strconv.ParseFloat(filter.Value, 32)

		if err != nil {
			return fmt.Errorf("%s: %w", filter.Value, ErrInvalidFilterValue)
		}		
	default:
		return fmt.Errorf("%s: %w", filter.Operator, ErrInvalidFilterOperator)
	}

	return nil
}

func validateSearchOrder(order SearchOrder) error {
	switch order.Field {
	case "name", "brand", "flavor", "carbohydrates", "electrolytes", "maltodextrose", "fructose", "caffeine", "sodium", "protein":
		switch order.Direction {
		case "asc", "desc":
			return nil
		default:
			return fmt.Errorf("%s: %w", order.Direction, ErrInvalidOrderDirection)
		}
	default:
		return fmt.Errorf("%s: %w", order.Field, ErrInvalidOrderField)
	}
}

func validateLimit(limit int) error {
	const maxLimit = 100

	if limit < 0 {
		return fmt.Errorf("%d: %w", limit, ErrInvalidLimit)
	}

	if limit > maxLimit {
		return fmt.Errorf("%d: %w", limit, ErrInvalidLimit)
	}

	return nil
}

func validateCursorDirection(cursor SearchCursor) error {
	switch cursor.Direction {
	case "next", "previous":
		return nil
	default:
		return fmt.Errorf("%s: %w", cursor.Direction, ErrInvalidCursorDirection)
	}
}

