package supplement

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type Supplement struct {
	Gtin          string  `json:"gtin"`
	Name          string  `json:"name"`
	Brand         string  `json:"brand"`
	Flavor        string  `json:"flavor"`
	Carbohydrates float32 `json:"carbohydrates"`
	Electrolytes  float32 `json:"electrolytes"`
	Maltodextrose float32 `json:"maltodextrose"`
	Fructose      float32 `json:"fructose"`
	Caffeine      float32 `json:"caffeine"`
	Sodium        float32 `json:"sodium"`
	Protein       float32 `json:"protein"`
}

type UpdatableSupplement struct {
	Name          *string  `json:"name,omitempty"`
	Brand         *string  `json:"brand,omitempty"`
	Flavor        *string  `json:"flavor,omitempty"`
	Carbohydrates *float32 `json:"carbohydrates,omitempty"`
	Electrolytes  *float32 `json:"electrolytes,omitempty"`
	Maltodextrose *float32 `json:"maltodextrose,omitempty"`
	Fructose      *float32 `json:"fructose,omitempty"`
	Caffeine      *float32 `json:"caffeine,omitempty"`
	Sodium        *float32 `json:"sodium,omitempty"`
	Protein       *float32 `json:"protein,omitempty"`
}

type SupplementRepository interface {
	FindByGtin(ctx context.Context, gtin string) (*Supplement, error)
	Create(ctx context.Context, supplement Supplement) error
	Update(ctx context.Context, supplement Supplement) error
	Delete(ctx context.Context, supplement Supplement) error
	ListAll(ctx context.Context) ([]Supplement, error)
}

func (s *Supplement) validate() error {
	var errors []string

	if !regexp.MustCompile(`^\d{13}$`).MatchString(s.Gtin) {
		errors = append(errors, fmt.Sprintf("gtin %q is invalid, it must be a 13-digit number", s.Gtin))
	}

	if s.Name == "" {
		errors = append(errors, fmt.Sprintf("name %q is invalid, it must not be empty", s.Name))
	}

	if s.Brand == "" {
		errors = append(errors, fmt.Sprintf("brand %q is invalid, it must not be empty", s.Brand))
	}

	if s.Flavor == "" {
		errors = append(errors, fmt.Sprintf("flavor %q is invalid, it must not be empty", s.Flavor))
	}

	if s.Carbohydrates < 0 {
		errors = append(errors, fmt.Sprintf("carbohydrates %f is invalid, it must be greater or equal to zero", s.Carbohydrates))
	}

	if s.Electrolytes < 0 {
		errors = append(errors, fmt.Sprintf("electrolytes %f is invalid, it must be greater or equal to zero", s.Electrolytes))
	}

	if s.Maltodextrose < 0 {
		errors = append(errors, fmt.Sprintf("maltodextrose %f is invalid, it must be greater or equal to zero", s.Maltodextrose))
	}

	if s.Fructose < 0 {
		errors = append(errors, fmt.Sprintf("fructose %f is invalid, it must be greater or equal to zero", s.Fructose))
	}

	if s.Caffeine < 0 {
		errors = append(errors, fmt.Sprintf("caffeine %f is invalid, it must be greater or equal to zero", s.Caffeine))
	}

	if s.Sodium < 0 {
		errors = append(errors, fmt.Sprintf("sodium %f is invalid, it must be greater or equal to zero", s.Sodium))
	}

	if s.Protein < 0 {
		errors = append(errors, fmt.Sprintf("protein %f is invalid, it must be greater or equal to zero", s.Protein))
	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(errors, "; "))
}

func (s *Supplement) update(other UpdatableSupplement) Supplement {
	if other.Name != nil {
		s.Name = *other.Name
	}

	if other.Brand != nil {
		s.Brand = *other.Brand
	}

	if other.Flavor != nil {
		s.Flavor = *other.Flavor
	}

	if other.Carbohydrates != nil {
		s.Carbohydrates = *other.Carbohydrates
	}

	if other.Electrolytes != nil {
		s.Electrolytes = *other.Electrolytes
	}

	if other.Maltodextrose != nil {
		s.Maltodextrose = *other.Maltodextrose
	}

	if other.Fructose != nil {
		s.Fructose = *other.Fructose
	}

	if other.Caffeine != nil {
		s.Caffeine = *other.Caffeine
	}

	if other.Sodium != nil {
		s.Sodium = *other.Sodium
	}

	if other.Protein != nil {
		s.Protein = *other.Protein
	}

	return *s
}
