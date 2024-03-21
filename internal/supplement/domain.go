package supplement

import (
	"context"
)

type Supplement struct {
	Gtin          string `json:"gtin"`
	Name          string `json:"name"`
	Brand         string `json:"brand"`
	Flavor        string `json:"flavor"`
	Carbohydrates float32    `json:"carbohydrates"`
	Electrolytes  float32    `json:"electrolytes"`
	Maltodextrose float32    `json:"maltodextrose"`
	Fructose      float32    `json:"fructose"`
	Caffeine      float32    `json:"caffeine"`
	Sodium        float32    `json:"sodium"`
	Protein       float32    `json:"protein"`
}

type UpdatableSupplement struct {
	Name          *string `json:"name,omitempty"`
	Brand         *string `json:"brand,omitempty"`
	Flavor        *string `json:"flavor,omitempty"`
	Carbohydrates *float32    `json:"carbohydrates,omitempty"`
	Electrolytes  *float32    `json:"electrolytes,omitempty"`
	Maltodextrose *float32    `json:"maltodextrose,omitempty"`
	Fructose      *float32    `json:"fructose,omitempty"`
	Caffeine      *float32    `json:"caffeine,omitempty"`
	Sodium        *float32    `json:"sodium,omitempty"`
	Protein       *float32    `json:"protein,omitempty"`
}

type SupplementRepository interface {
	FindByGtin(ctx context.Context, gtin string) (*Supplement, error)
	Create(ctx context.Context, supplement Supplement) error
	Update(ctx context.Context, supplement Supplement) error
	Delete(ctx context.Context, supplement Supplement) error
	ListAll(ctx context.Context) ([]Supplement, error)
}

func (supplement *Supplement) update(other UpdatableSupplement) Supplement {
	if other.Name != nil {
		supplement.Name = *other.Name
	}

	if other.Brand != nil {
		supplement.Brand = *other.Brand
	}

	if other.Flavor != nil {
		supplement.Flavor = *other.Flavor
	}

	if other.Carbohydrates != nil {
		supplement.Carbohydrates = *other.Carbohydrates
	}

	if other.Electrolytes != nil {
		supplement.Electrolytes = *other.Electrolytes
	}

	if other.Maltodextrose != nil {
		supplement.Maltodextrose = *other.Maltodextrose
	}

	if other.Fructose != nil {
		supplement.Fructose = *other.Fructose
	}

	if other.Caffeine != nil {
		supplement.Caffeine = *other.Caffeine
	}

	if other.Sodium != nil {
		supplement.Sodium = *other.Sodium
	}

	if other.Protein != nil {
		supplement.Protein = *other.Protein
	}

	return *supplement
}
