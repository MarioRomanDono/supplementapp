package supplement

type Supplement struct {
	Gtin          string `json:"gtin"`
	Name          string `json:"name"`
	Brand         string `json:"brand"`
	Flavor        string `json:"flavor"`
	Carbohydrates int    `json:"carbohydrates"`
	Electrolytes  int    `json:"electrolytes"`
	Maltodextrose int    `json:"maltodextrose"`
	Fructose      int    `json:"fructose"`
	Caffeine      int    `json:"caffeine"`
	Sodium        int    `json:"sodium"`
	Protein       int    `json:"protein"`
}

type UpdatableSupplement struct {
	Name          *string `json:"name,omitempty"`
	Brand         *string `json:"brand,omitempty"`
	Flavor        *string `json:"flavor,omitempty"`
	Carbohydrates *int    `json:"carbohydrates,omitempty"`
	Electrolytes  *int    `json:"electrolytes,omitempty"`
	Maltodextrose *int    `json:"maltodextrose,omitempty"`
	Fructose      *int    `json:"fructose,omitempty"`
	Caffeine      *int    `json:"caffeine,omitempty"`
	Sodium        *int    `json:"sodium,omitempty"`
	Protein       *int    `json:"protein,omitempty"`
}

type SupplementRepository interface {
	FindByGtin(gtin string) (*Supplement, error)
	Create(supplement Supplement) error
	Update(supplement Supplement) error
	Delete(supplement Supplement) error
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
