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

func TestCreate_NonExisting(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	randomSupplement := newRandomSupplement()

	err := service.Create(*randomSupplement)

	if err != nil {
		t.Errorf("expected err to be nil; got: %s", err)
	}

	if len(repository.createCalls) != 1 {
		t.Errorf("expected create to be called just once; called: %d", len(repository.createCalls))
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
		t.Errorf("expected create to be called just once; called: %d", len(repository.createCalls))
	}

	if len(repository.findByGtinCalls) != 1 {
		t.Errorf("expected findByGtinCalls to be called just once; called: %d", len(repository.findByGtinCalls))
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
		t.Errorf("expected findByGtinCalls to be called just once; called: %d", len(repository.findByGtinCalls))
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
		t.Errorf("expected findByGtinCalls to be called just once; called: %d", len(repository.findByGtinCalls))
	}

	lastFoundGtin := repository.findByGtinCalls[0]

	if lastFoundGtin != randomSupplement.Gtin {
		t.Errorf("expected %s; got %s", randomSupplement.Gtin, lastFoundGtin)
	}
}
