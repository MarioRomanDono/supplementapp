package supplement_test

import (
	"errors"
	"strconv"
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
	searchCalls	    []struct {
		Filters []supplement.SearchFilter
		Order   *supplement.SearchOrder
		Limit   int
		Cursor  *supplement.SearchCursor
	}
	searchResponse *supplement.SearchResponse
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

func (repository *mockSupplementRepository) Search(
	filters []supplement.SearchFilter,
	order *supplement.SearchOrder,
	limit int,
	cursor *supplement.SearchCursor) (*supplement.SearchResponse, error) {
	repository.searchCalls = append(repository.searchCalls, struct {
		Filters []supplement.SearchFilter
		Order   *supplement.SearchOrder
		Limit   int
		Cursor  *supplement.SearchCursor
	}{filters, order, limit, cursor})

	return repository.searchResponse, nil
}

func (repository *mockSupplementRepository) whenSearchThenReturn(response *supplement.SearchResponse) {
	repository.searchResponse = response
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
		Carbohydrates: gofakeit.Float32Range(0, 100),
		Electrolytes:  gofakeit.Float32Range(0, 100),
		Maltodextrose: gofakeit.Float32Range(0, 100),
		Fructose:      gofakeit.Float32Range(0, 100),
		Caffeine:      gofakeit.Float32Range(0, 100),
		Sodium:        gofakeit.Float32Range(0, 100),
		Protein:       gofakeit.Float32Range(0, 100),
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
	var carbohydrates *float32
	var electrolytes *float32
	var maltodextrose *float32
	var fructose *float32
	var caffeine *float32
	var sodium *float32
	var protein *float32

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
		c := gofakeit.Float32Range(0, 100)
		carbohydrates = &c
	}

	if gofakeit.Bool() {
		e := gofakeit.Float32Range(0, 100)
		electrolytes = &e
	}

	if gofakeit.Bool() {
		m := gofakeit.Float32Range(0, 100)
		maltodextrose = &m
	}

	if gofakeit.Bool() {
		f := gofakeit.Float32Range(0, 100)
		fructose = &f
	}

	if gofakeit.Bool() {
		c := gofakeit.Float32Range(0, 100)
		caffeine = &c
	}

	if gofakeit.Bool() {
		s := gofakeit.Float32Range(0, 100)
		sodium = &s
	}

	if gofakeit.Bool() {
		p := gofakeit.Float32Range(0, 100)
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

func generateRandomSearchOptions() supplement.SearchOptions {
	var filters []supplement.SearchFilter
	var order *supplement.SearchOrder
	var limit int
	var cursor *supplement.SearchCursor

	if gofakeit.Bool() {
		size := gofakeit.Number(1, 10)
		
		for i := 0; i < size; i++ {
			if gofakeit.Bool() {
				filters = append(filters, generateRandomStringFilter())
			} else {
				filters = append(filters, generateRandomFloatFilter())
			}
		}
	}

	if gofakeit.Bool() {
		order = &supplement.SearchOrder{
			Field: gofakeit.RandomString([]string{"name", "brand", "flavor", "carbohydrates", "electrolytes", "maltodextrose", "fructose", "caffeine", "sodium", "protein"}),
			Direction: gofakeit.RandomString([]string{"asc", "desc"}),
		}
	}

	if gofakeit.Bool() {
		limit = gofakeit.Number(1, 100)
	}

	if gofakeit.Bool() {
		cursor = &supplement.SearchCursor{
			Cursor: gofakeit.Word(),
			Direction: gofakeit.RandomString([]string{"next", "previous"}),
		}
	}

	return supplement.SearchOptions{
		Filters: filters,
		Order: order,
		Limit: limit,
		Cursor: cursor,
	}
}

func generateRandomStringFilter() supplement.SearchFilter {
	return supplement.SearchFilter{
		Field: gofakeit.RandomString([]string{"name", "brand", "flavor"}),
		Operator: gofakeit.RandomString([]string{"=", "!=", "contains"}),
		Value: gofakeit.Word(),
	}
}

func generateRandomFloatFilter() supplement.SearchFilter {
	return supplement.SearchFilter{
		Field: gofakeit.RandomString([]string{"carbohydrates", "electrolytes", "maltodextrose", "fructose", "caffeine", "sodium", "protein"}),
		Operator: gofakeit.RandomString([]string{"=", "!=", ">", "<", ">=", "<="}),
		Value: strconv.FormatFloat(gofakeit.Float64(), 'f', -1, 32),
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

func TestSearch_InvalidFilterFields(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	filters := []supplement.SearchFilter{
		{
			Field: "invalid",
		},
	}
	options := supplement.SearchOptions{
		Filters: filters,
	}

	_, err := service.Search(options)

	if !errors.Is(err, supplement.ErrInvalidFilterField) {
		t.Errorf("expected err to be ErrInvalidFilterField; got: %s", err)
	}
}

func TestSearch_InvalidFilterOperator(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	filters := []supplement.SearchFilter{
		{
			Field: gofakeit.RandomString([]string{"name", "brand", "flavor", "carbohydrates", "electrolytes", "maltodextrose", "fructose", "caffeine", "sodium", "protein"}),
			Operator: "invalid",
		},
	}
	options := supplement.SearchOptions{
		Filters: filters,
	}

	_, err := service.Search(options)

	if !errors.Is(err, supplement.ErrInvalidFilterOperator) {
		t.Errorf("expected err to be ErrInvalidFilterOperator; got: %s", err)
	}
}

func TestSearch_InvalidFilterValueFloat(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	filters := []supplement.SearchFilter{
		{
			Field: gofakeit.RandomString([]string{"carbohydrates", "electrolytes", "maltodextrose", "fructose", "caffeine", "sodium", "protein"}),
			Operator: gofakeit.RandomString([]string{"=", "!=", ">", "<", ">=", "<="}),
			Value: "invalid",
		},
	}
	options := supplement.SearchOptions{
		Filters: filters,
	}

	_, err := service.Search(options)

	if !errors.Is(err, supplement.ErrInvalidFilterValue) {
		t.Errorf("expected err to be ErrInvalidFilterValue; got: %s", err)
	}
}

func TestSearch_InvalidOrderField(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	order := supplement.SearchOrder{
		Field: "invalid",
	}
	options := supplement.SearchOptions{
		Order: &order,
	}

	_, err := service.Search(options)

	if !errors.Is(err, supplement.ErrInvalidOrderField) {
		t.Errorf("expected err to be ErrInvalidOrderField; got: %s", err)
	}
}

func TestSearch_InvalidOrderDirection(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	order := supplement.SearchOrder{
		Field: gofakeit.RandomString([]string{"name", "brand", "flavor", "carbohydrates", "electrolytes", "maltodextrose", "fructose", "caffeine", "sodium", "protein"}),
		Direction: "invalid",
	}
	options := supplement.SearchOptions{
		Order: &order,
	}

	_, err := service.Search(options)

	if !errors.Is(err, supplement.ErrInvalidOrderDirection) {
		t.Errorf("expected err to be ErrInvalidOrderDirection; got: %s", err)
	}
}

func TestSearch_InvalidLimit(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	limit := gofakeit.Number(-100, -1)
	options := supplement.SearchOptions{
		Limit: limit,
	}

	_, err := service.Search(options)

	if !errors.Is(err, supplement.ErrInvalidLimit) {
		t.Errorf("expected err to be ErrInvalidLimit; got: %s", err)
	}
}

func TestSearch_InvalidCursorDirection(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	cursor := supplement.SearchCursor{
		Direction: "invalid",
	}
	options := supplement.SearchOptions{
		Cursor: &cursor,
	}

	_, err := service.Search(options)

	if !errors.Is(err, supplement.ErrInvalidCursorDirection) {
		t.Errorf("expected err to be ErrInvalidCursorDirection; got: %s", err)
	}
}

func TestSearch(t *testing.T) {
	repository := newMockSupplementRepository()
	service := supplement.NewSupplementService(repository)
	options := generateRandomSearchOptions()
	
	var expected *supplement.SearchResponse
	gofakeit.Struct(&expected)
	repository.whenSearchThenReturn(expected)

	response, err := service.Search(options)

	if err != nil {
		t.Errorf("expected err to be nil; got: %s", err)
	}

	if len(repository.searchCalls) != 1 {
		t.Errorf("expected Search to be called just once; called: %d", len(repository.searchCalls))
	}

	lastSearch := repository.searchCalls[0]
	searchArgs := struct {
		Filters []supplement.SearchFilter
		Order   *supplement.SearchOrder
		Limit   int
		Cursor  *supplement.SearchCursor
	}{options.Filters, options.Order, options.Limit, options.Cursor}

	if !cmp.Equal(lastSearch, searchArgs) {
		t.Errorf("expected %v; got %v", searchArgs, lastSearch)
	}

	if !cmp.Equal(response, expected) {
		t.Errorf("expected %v; got %v", expected, response)
	}
}