package supplement_test

import (
	"errors"
	"testing"

	"github.com/marioromandono/supplementapp/internal/supplement"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/go-cmp/cmp"
)

type mockSupplementRepository struct {
	store           map[string]supplement.Supplement
	findByGtinCalls []string
	createCalls     []supplement.Supplement
	updateCalls     []supplement.Supplement
	deleteCalls     []supplement.Supplement
}

func (repository *mockSupplementRepository) FindByGtin(gtin string) (*supplement.Supplement, error) {
	repository.findByGtinCalls = append(repository.findByGtinCalls, gtin)
	stored, ok := repository.store[gtin]

	if !ok {
		return nil, nil
	}

	return &stored, nil
}

func (repository *mockSupplementRepository) Create(supplement supplement.Supplement) error {
	repository.store[supplement.Gtin] = supplement
	repository.createCalls = append(repository.createCalls, supplement)
	return nil
}

func (repository *mockSupplementRepository) Update(supplement supplement.Supplement) error {
	repository.store[supplement.Gtin] = supplement
	repository.updateCalls = append(repository.updateCalls, supplement)
	return nil
}

func (repository *mockSupplementRepository) Delete(supplement supplement.Supplement) error {
	delete(repository.store, supplement.Gtin)
	repository.deleteCalls = append(repository.deleteCalls, supplement)
	return nil
}

func newMockSupplementRepository() *mockSupplementRepository {
	return &mockSupplementRepository{
		store: make(map[string]supplement.Supplement),
	}
}

func newRandomSupplement() *supplement.Supplement {
	return &supplement.Supplement{
		Gtin:          gofakeit.DigitN(13),
		Name:          gofakeit.Name(),
		Brand:         gofakeit.Name(),
		Flavor:        gofakeit.Word(),
		Carbohydrates: gofakeit.Number(0, 100),
		Electrolytes:  gofakeit.Number(0, 100),
		Maltodextrose: gofakeit.Number(0, 100),
		Fructose:      gofakeit.Number(0, 100),
		Caffeine:      gofakeit.Number(0, 100),
		Sodium:        gofakeit.Number(0, 100),
		Protein:       gofakeit.Number(0, 100),
	}
}

func updateSupplement(supplement supplement.Supplement, other supplement.UpdatableSupplement) supplement.Supplement {
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

	return supplement
}

func newRandomUpdatableSupplement() *supplement.UpdatableSupplement {
	var name *string
	var brand *string
	var flavor *string
	var carbohydrates *int
	var electrolytes *int
	var maltodextrose *int
	var fructose *int
	var caffeine *int
	var sodium *int
	var protein *int

	if gofakeit.Bool() {
		n := gofakeit.Name()
		name = &n
	}

	if gofakeit.Bool() {
		b := gofakeit.Name()
		brand = &b
	}

	if gofakeit.Bool() {
		f := gofakeit.Word()
		flavor = &f
	}

	if gofakeit.Bool() {
		c := gofakeit.Number(0, 100)
		carbohydrates = &c
	}

	if gofakeit.Bool() {
		e := gofakeit.Number(0, 100)
		electrolytes = &e
	}

	if gofakeit.Bool() {
		m := gofakeit.Number(0, 100)
		maltodextrose = &m
	}

	if gofakeit.Bool() {
		f := gofakeit.Number(0, 100)
		fructose = &f
	}

	if gofakeit.Bool() {
		c := gofakeit.Number(0, 100)
		caffeine = &c
	}

	if gofakeit.Bool() {
		s := gofakeit.Number(0, 100)
		sodium = &s
	}

	if gofakeit.Bool() {
		p := gofakeit.Number(0, 100)
		protein = &p
	}

	return &supplement.UpdatableSupplement{
		Name:          name,
		Brand:         brand,
		Flavor:        flavor,
		Carbohydrates: carbohydrates,
		Electrolytes:  electrolytes,
		Maltodextrose: maltodextrose,
		Fructose:      fructose,
		Caffeine:      caffeine,
		Sodium:        sodium,
		Protein:       protein,
	}
}

func TestCreate_NonExisting(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	randomSupplement := newRandomSupplement()

	err := service.Create(*randomSupplement)

	if err != nil {
		t.Errorf("expected err to be nil; got: %s", err)
	}

	if len(repository.createCalls) != 1 {
		t.Errorf("expected Create to be called just once; called: %d", len(repository.createCalls))
	}

	lastCreatedSupplement := repository.createCalls[0]

	if lastCreatedSupplement != *randomSupplement {
		t.Errorf("expected %v; got %v", randomSupplement, lastCreatedSupplement)
	}
}

func TestCreate_Existing(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	randomSupplement := newRandomSupplement()

	err := repository.Create(*randomSupplement)

	if err != nil {
		t.Errorf("expected err to be nil; got: %s", err)
	}

	err = service.Create(*randomSupplement)

	if !errors.Is(err, supplement.ErrAlreadyExists) {
		t.Errorf("expected err to be ErrAlreadyExists; got: %s", err)
	}

	if len(repository.createCalls) != 1 {
		t.Errorf("expected Create to be called just once; called: %d", len(repository.createCalls))
	}

	if len(repository.findByGtinCalls) != 1 {
		t.Errorf("expected FindByGtin to be called just once; called: %d", len(repository.findByGtinCalls))
	}

	lastFoundGtin := repository.findByGtinCalls[0]

	if lastFoundGtin != randomSupplement.Gtin {
		t.Errorf("expected %s; got %s", randomSupplement.Gtin, lastFoundGtin)
	}
}

func TestFindByGtin_NonExisting(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	gtin := gofakeit.DigitN(13)

	_, err := service.FindByGtin(gtin)

	if !errors.Is(err, supplement.ErrNotFound) {
		t.Errorf("expected err to be ErrNotFound; got: %s", err)
	}

	if len(repository.findByGtinCalls) != 1 {
		t.Errorf("expected FindByGtin to be called just once; called: %d", len(repository.findByGtinCalls))
	}

	lastFoundGtin := repository.findByGtinCalls[0]

	if lastFoundGtin != gtin {
		t.Errorf("expected %s; got %s", gtin, lastFoundGtin)
	}
}

func TestFindByGtin_Existing(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	randomSupplement := newRandomSupplement()

	err := repository.Create(*randomSupplement)

	if err != nil {
		t.Errorf("expected err to be nil; got: %s", err)
	}

	foundSupplement, err := service.FindByGtin(randomSupplement.Gtin)

	if err != nil {
		t.Errorf("expected err to be nil; got: %s", err)
	}

	if !cmp.Equal(randomSupplement, foundSupplement) {
		t.Errorf("expected %v; got %v", randomSupplement, foundSupplement)
	}

	if len(repository.findByGtinCalls) != 1 {
		t.Errorf("expected FindByGtin to be called just once; called: %d", len(repository.findByGtinCalls))
	}

	lastFoundGtin := repository.findByGtinCalls[0]

	if lastFoundGtin != randomSupplement.Gtin {
		t.Errorf("expected %s; got %s", randomSupplement.Gtin, lastFoundGtin)
	}
}

func TestDelete_NonExisting(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	gtin := gofakeit.DigitN(13)

	err := service.Delete(gtin)

	if !errors.Is(err, supplement.ErrNotFound) {
		t.Errorf("expected err to be ErrNotFound; got: %s", err)
	}

	if len(repository.findByGtinCalls) != 1 {
		t.Errorf("expected FindByGtin to be called just once; called: %d", len(repository.findByGtinCalls))
	}

	lastFoundGtin := repository.findByGtinCalls[0]

	if lastFoundGtin != gtin {
		t.Errorf("expected %s; got %s", gtin, lastFoundGtin)
	}
}

func TestDelete_Existing(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	randomSupplement := newRandomSupplement()

	err := repository.Create(*randomSupplement)

	if err != nil {
		t.Errorf("expected err to be nil; got: %s", err)
	}

	err = service.Delete(randomSupplement.Gtin)

	if err != nil {
		t.Errorf("expected err to be nil; got: %s", err)
	}

	if len(repository.findByGtinCalls) != 1 {
		t.Errorf("expected FindByGtin to be called just once; called: %d", len(repository.findByGtinCalls))
	}

	lastFoundGtin := repository.findByGtinCalls[0]

	if lastFoundGtin != randomSupplement.Gtin {
		t.Errorf("expected %s; got %s", randomSupplement.Gtin, lastFoundGtin)
	}

	if len(repository.deleteCalls) != 1 {
		t.Errorf("expected Delete to be called just once; called: %d", len(repository.deleteCalls))
	}

	lastDeleted := repository.deleteCalls[0]

	if !cmp.Equal(lastDeleted, *randomSupplement) {
		t.Errorf("expected %v; got %v", randomSupplement, lastDeleted)
	}
}

func TestUpdate_NonExisting(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	gtin := gofakeit.DigitN(13)
	updatableSupplement := newRandomUpdatableSupplement()

	err := service.Update(gtin, *updatableSupplement)

	if !errors.Is(err, supplement.ErrNotFound) {
		t.Errorf("expected err to be ErrNotFound; got: %s", err)
	}

	if len(repository.findByGtinCalls) != 1 {
		t.Errorf("expected FindByGtin to be called just once; called: %d", len(repository.findByGtinCalls))
	}

	lastFoundGtin := repository.findByGtinCalls[0]

	if lastFoundGtin != gtin {
		t.Errorf("expected %s; got %s", gtin, lastFoundGtin)
	}
}

func TestUpdate_Existing(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	randomSupplement := newRandomSupplement()

	err := repository.Create(*randomSupplement)

	if err != nil {
		t.Errorf("expected err to be nil; got: %s", err)
	}

	updatableSupplement := newRandomUpdatableSupplement()

	err = service.Update(randomSupplement.Gtin, *updatableSupplement)

	if err != nil {
		t.Errorf("expected err to be nil; got: %s", err)
	}

	if len(repository.findByGtinCalls) != 1 {
		t.Errorf("expected FindByGtin to be called just once; called: %d", len(repository.findByGtinCalls))
	}

	lastFoundGtin := repository.findByGtinCalls[0]

	if lastFoundGtin != randomSupplement.Gtin {
		t.Errorf("expected %s; got %s", randomSupplement.Gtin, lastFoundGtin)
	}

	if len(repository.updateCalls) != 1 {
		t.Errorf("expected Update to be called just once; called: %d", len(repository.updateCalls))
	}

	lastUpdated := repository.updateCalls[0]
	expectedUpdated := updateSupplement(*randomSupplement, *updatableSupplement)

	if !cmp.Equal(lastUpdated, expectedUpdated) {
		t.Errorf("expected %v; got %v", expectedUpdated, lastUpdated)
	}
}
