package supplement

type Supplement struct {
	Gtin string `json:"gtin"`
	Name string `json:"name"`
	Brand string `json:"brand"`
	Flavor string `json:"flavor"`
	Carbohydrates int `json:"carbohydrates"`
	Electrolytes int `json:"electrolytes"`
	Maltodextrose int `json:"maltodextrose"`
	Fructose int `json:"fructose"`
	Caffeine int `json:"caffeine"`
	Sodium int `json:"sodium"`
	Protein int `json:"protein"`
}

type SupplementRepository interface {
	FindByGtin(gtin string) (*Supplement, error)
	Create(supplement Supplement) error
	Update(supplement Supplement) error
	Delete(supplement Supplement) error
}