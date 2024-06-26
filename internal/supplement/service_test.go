package supplement_test

import (
	"context"
	"errors"
	"testing"

	"github.com/marioromandono/supplementapp/internal/supplement"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type stubSupplementRepository struct {
	store map[string]supplement.Supplement
}

func (r *stubSupplementRepository) FindByGtin(ctx context.Context, gtin string) (*supplement.Supplement, error) {
	s, ok := r.store[gtin]
	if !ok {
		return nil, nil
	}
	return &s, nil
}

func (r *stubSupplementRepository) Create(ctx context.Context, s supplement.Supplement) error {
	r.store[s.Gtin] = s
	return nil
}

func (r *stubSupplementRepository) Update(ctx context.Context, s supplement.Supplement) error {
	r.store[s.Gtin] = s
	return nil
}

func (r *stubSupplementRepository) Delete(ctx context.Context, s supplement.Supplement) error {
	delete(r.store, s.Gtin)
	return nil
}

func (r *stubSupplementRepository) ListAll(ctx context.Context) ([]supplement.Supplement, error) {
	var supplements []supplement.Supplement
	for _, s := range r.store {
		supplements = append(supplements, s)
	}
	return supplements, nil
}

func Ptr[T any](v T) *T {
	return &v
}

func TestSupplementService_FindByGtin(t *testing.T) {
	t.Parallel()
	type fields struct {
		repository supplement.SupplementRepository
	}
	type args struct {
		ctx  context.Context
		gtin string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *supplement.Supplement
		wantErr   error
		wantStore map[string]supplement.Supplement
	}{
		{
			name: "not found",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx:  context.TODO(),
				gtin: "1234567890123",
			},
			want:      nil,
			wantErr:   supplement.ErrNotFound,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "found",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {Gtin: "1234567890123"},
				}},
			},
			args: args{
				ctx:  context.TODO(),
				gtin: "1234567890123",
			},
			want:    &supplement.Supplement{Gtin: "1234567890123"},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {Gtin: "1234567890123"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := supplement.NewSupplementService(tt.fields.repository)
			got, err := service.FindByGtin(tt.args.ctx, tt.args.gtin)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SupplementService.FindByGtin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("SupplementService.FindByGtin() = %v, want %v", got, tt.want)
			}
			if diff := cmp.Diff(tt.fields.repository.(*stubSupplementRepository).store, tt.wantStore); diff != "" {
				t.Errorf("SupplementService.FindByGtin() store mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestSupplementService_Create(t *testing.T) {
	t.Parallel()
	type fields struct {
		repository supplement.SupplementRepository
	}
	type args struct {
		ctx        context.Context
		supplement supplement.Supplement
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   error
		wantStore map[string]supplement.Supplement
	}{
		{
			name: "already exists",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {Gtin: "1234567890123"},
				}},
			},
			args: args{
				ctx:        context.TODO(),
				supplement: supplement.Supplement{Gtin: "1234567890123"},
			},
			wantErr: supplement.ErrAlreadyExists,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {Gtin: "1234567890123"},
			},
		},
		{
			name: "missing gtin",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx:        context.TODO(),
				supplement: supplement.Supplement{},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "invalid gtin",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx:        context.TODO(),
				supplement: supplement.Supplement{Gtin: "1"},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "missing name",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx:        context.TODO(),
				supplement: supplement.Supplement{Gtin: "1234567890123"},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "missing brand",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin: "1234567890123",
					Name: "name",
				},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "missing flavor",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin:  "1234567890123",
					Name:  "name",
					Brand: "brand",
				},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "invalid carbohydrates",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin:          "1234567890123",
					Name:          "name",
					Brand:         "brand",
					Flavor:        "flavor",
					Carbohydrates: -1.0,
				},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "invalid electrolytes",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin:         "1234567890123",
					Name:         "name",
					Brand:        "brand",
					Flavor:       "flavor",
					Electrolytes: -1.0,
				},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "invalid maltodextrose",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin:          "1234567890123",
					Name:          "name",
					Brand:         "brand",
					Flavor:        "flavor",
					Maltodextrose: -1.0,
				},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "invalid fructose",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin:     "1234567890123",
					Name:     "name",
					Brand:    "brand",
					Flavor:   "flavor",
					Fructose: -1.0,
				},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "invalid caffeine",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin:     "1234567890123",
					Name:     "name",
					Brand:    "brand",
					Flavor:   "flavor",
					Caffeine: -1.0,
				},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "invalid sodium",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
					Sodium: -1.0,
				},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "invalid protein",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin:    "1234567890123",
					Name:    "name",
					Brand:   "brand",
					Flavor:  "flavor",
					Protein: -1.0,
				},
			},
			wantErr:   supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "with every nutrient set to zero",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "with every nutrient set to a positive value",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
				supplement: supplement.Supplement{
					Gtin:          "1234567890123",
					Name:          "name",
					Brand:         "brand",
					Flavor:        "flavor",
					Carbohydrates: 1.0,
					Electrolytes:  1.0,
					Maltodextrose: 1.0,
					Fructose:      1.0,
					Caffeine:      1.0,
					Sodium:        1.0,
					Protein:       1.0,
				},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:          "1234567890123",
					Name:          "name",
					Brand:         "brand",
					Flavor:        "flavor",
					Carbohydrates: 1.0,
					Electrolytes:  1.0,
					Maltodextrose: 1.0,
					Fructose:      1.0,
					Caffeine:      1.0,
					Sodium:        1.0,
					Protein:       1.0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := supplement.NewSupplementService(tt.fields.repository)
			if err := service.Create(tt.args.ctx, tt.args.supplement); !errors.Is(err, tt.wantErr) {
				t.Errorf("SupplementService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.fields.repository.(*stubSupplementRepository).store, tt.wantStore); diff != "" {
				t.Errorf("SupplementService.Create() store mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestSupplementService_Update(t *testing.T) {
	t.Parallel()
	type fields struct {
		repository supplement.SupplementRepository
	}
	type args struct {
		ctx   context.Context
		gtin  string
		other supplement.UpdatableSupplement
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   error
		wantStore map[string]supplement.Supplement
	}{
		{
			name: "not found",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Name: Ptr("updated name")},
			},
			wantErr:   supplement.ErrNotFound,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "invalid name",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Name: Ptr("")},
			},
			wantErr: supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "valid name",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Name: Ptr("updated name")},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "updated name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "invalid brand",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Brand: Ptr("")},
			},
			wantErr: supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "valid brand",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Brand: Ptr("updated brand")},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "updated brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "invalid flavor",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Flavor: Ptr("")},
			},
			wantErr: supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "valid flavor",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Flavor: Ptr("updated flavor")},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "updated flavor",
				},
			},
		},
		{
			name: "invalid carbohydrates",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Carbohydrates: Ptr(float32(-1.0))},
			},
			wantErr: supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "valid carbohydrates",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:          "1234567890123",
						Name:          "name",
						Brand:         "brand",
						Flavor:        "flavor",
						Carbohydrates: 1.0,
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Carbohydrates: Ptr(float32(2.0))},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:          "1234567890123",
					Name:          "name",
					Brand:         "brand",
					Flavor:        "flavor",
					Carbohydrates: 2.0,
				},
			},
		},
		{
			name: "invalid electrolytes",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Electrolytes: Ptr(float32(-1.0))},
			},
			wantErr: supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "valid electrolytes",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:         "1234567890123",
						Name:         "name",
						Brand:        "brand",
						Flavor:       "flavor",
						Electrolytes: 1.0,
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Electrolytes: Ptr(float32(2.0))},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:         "1234567890123",
					Name:         "name",
					Brand:        "brand",
					Flavor:       "flavor",
					Electrolytes: 2.0,
				},
			},
		},
		{
			name: "invalid maltodextrose",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Maltodextrose: Ptr(float32(-1.0))},
			},
			wantErr: supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "valid maltodextrose",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:          "1234567890123",
						Name:          "name",
						Brand:         "brand",
						Flavor:        "flavor",
						Maltodextrose: 1.0,
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Maltodextrose: Ptr(float32(2.0))},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:          "1234567890123",
					Name:          "name",
					Brand:         "brand",
					Flavor:        "flavor",
					Maltodextrose: 2.0,
				},
			},
		},
		{
			name: "invalid fructose",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Fructose: Ptr(float32(-1.0))},
			},
			wantErr: supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "valid fructose",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:     "1234567890123",
						Name:     "name",
						Brand:    "brand",
						Flavor:   "flavor",
						Fructose: 1.0,
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Fructose: Ptr(float32(2.0))},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:     "1234567890123",
					Name:     "name",
					Brand:    "brand",
					Flavor:   "flavor",
					Fructose: 2.0,
				},
			},
		},
		{
			name: "invalid caffeine",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Caffeine: Ptr(float32(-1.0))},
			},
			wantErr: supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "valid caffeine",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:     "1234567890123",
						Name:     "name",
						Brand:    "brand",
						Flavor:   "flavor",
						Caffeine: 1.0,
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Caffeine: Ptr(float32(2.0))},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:     "1234567890123",
					Name:     "name",
					Brand:    "brand",
					Flavor:   "flavor",
					Caffeine: 2.0,
				},
			},
		},
		{
			name: "invalid sodium",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Sodium: Ptr(float32(-1.0))},
			},
			wantErr: supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "valid sodium",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
						Sodium: 1.0,
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Sodium: Ptr(float32(2.0))},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
					Sodium: 2.0,
				},
			},
		},
		{
			name: "invalid protein",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:   "1234567890123",
						Name:   "name",
						Brand:  "brand",
						Flavor: "flavor",
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Protein: Ptr(float32(-1.0))},
			},
			wantErr: supplement.ErrInvalidSupplement,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:   "1234567890123",
					Name:   "name",
					Brand:  "brand",
					Flavor: "flavor",
				},
			},
		},
		{
			name: "valid protein",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:    "1234567890123",
						Name:    "name",
						Brand:   "brand",
						Flavor:  "flavor",
						Protein: 1.0,
					},
				}},
			},
			args: args{
				ctx:   context.TODO(),
				gtin:  "1234567890123",
				other: supplement.UpdatableSupplement{Protein: Ptr(float32(2.0))},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:    "1234567890123",
					Name:    "name",
					Brand:   "brand",
					Flavor:  "flavor",
					Protein: 2.0,
				},
			},
		},
		{
			name: "valid supplement",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {
						Gtin:          "1234567890123",
						Name:          "name",
						Brand:         "brand",
						Flavor:        "flavor",
						Carbohydrates: 1.0,
						Electrolytes:  1.0,
						Maltodextrose: 1.0,
						Fructose:      1.0,
						Caffeine:      1.0,
						Sodium:        1.0,
						Protein:       1.0,
					},
				}},
			},
			args: args{
				ctx:  context.TODO(),
				gtin: "1234567890123",
				other: supplement.UpdatableSupplement{
					Name:          Ptr("updated name"),
					Brand:         Ptr("updated brand"),
					Flavor:        Ptr("updated flavor"),
					Carbohydrates: Ptr(float32(2.0)),
					Electrolytes:  Ptr(float32(2.0)),
					Maltodextrose: Ptr(float32(2.0)),
					Fructose:      Ptr(float32(2.0)),
					Caffeine:      Ptr(float32(2.0)),
					Sodium:        Ptr(float32(2.0)),
					Protein:       Ptr(float32(2.0)),
				},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {
					Gtin:          "1234567890123",
					Name:          "updated name",
					Brand:         "updated brand",
					Flavor:        "updated flavor",
					Carbohydrates: 2.0,
					Electrolytes:  2.0,
					Maltodextrose: 2.0,
					Fructose:      2.0,
					Caffeine:      2.0,
					Sodium:        2.0,
					Protein:       2.0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := supplement.NewSupplementService(tt.fields.repository)
			if err := service.Update(tt.args.ctx, tt.args.gtin, tt.args.other); !errors.Is(err, tt.wantErr) {
				t.Errorf("SupplementService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.fields.repository.(*stubSupplementRepository).store, tt.wantStore); diff != "" {
				t.Errorf("SupplementService.Update() store mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestSupplementService_Delete(t *testing.T) {
	t.Parallel()
	type fields struct {
		repository supplement.SupplementRepository
	}
	type args struct {
		ctx  context.Context
		gtin string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   error
		wantStore map[string]supplement.Supplement
	}{
		{
			name: "not found",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx:  context.TODO(),
				gtin: "1234567890123",
			},
			wantErr:   supplement.ErrNotFound,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "success",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {Gtin: "1234567890123"},
				}},
			},
			args: args{
				ctx:  context.TODO(),
				gtin: "1234567890123",
			},
			wantErr:   nil,
			wantStore: map[string]supplement.Supplement{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := supplement.NewSupplementService(tt.fields.repository)
			if err := service.Delete(tt.args.ctx, tt.args.gtin); !errors.Is(err, tt.wantErr) {
				t.Errorf("SupplementService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.fields.repository.(*stubSupplementRepository).store, tt.wantStore); diff != "" {
				t.Errorf("SupplementService.Delete() store mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestSupplementService_ListAll(t *testing.T) {
	t.Parallel()
	type fields struct {
		repository supplement.SupplementRepository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args
		want      []supplement.Supplement
		wantErr   error
		wantStore map[string]supplement.Supplement
	}{
		{
			name: "empty store",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{}},
			},
			args: args{
				ctx: context.TODO(),
			},
			want:      []supplement.Supplement{},
			wantErr:   nil,
			wantStore: map[string]supplement.Supplement{},
		},
		{
			name: "non-empty store",
			fields: fields{
				repository: &stubSupplementRepository{store: map[string]supplement.Supplement{
					"1234567890123": {Gtin: "1234567890123"},
					"1234567890124": {Gtin: "1234567890124"},
				}},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: []supplement.Supplement{
				{Gtin: "1234567890123"},
				{Gtin: "1234567890124"},
			},
			wantErr: nil,
			wantStore: map[string]supplement.Supplement{
				"1234567890123": {Gtin: "1234567890123"},
				"1234567890124": {Gtin: "1234567890124"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := supplement.NewSupplementService(tt.fields.repository)
			got, err := service.ListAll(tt.args.ctx)
			less := func(a, b supplement.Supplement) bool {
				return a.Gtin < b.Gtin
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SupplementService.ListAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.EquateEmpty(), cmpopts.SortSlices(less)); diff != "" {
				t.Errorf("SupplementService.ListAll() (-got +want):\n%s", diff)
			}
			if diff := cmp.Diff(tt.fields.repository.(*stubSupplementRepository).store, tt.wantStore); diff != "" {
				t.Errorf("SupplementService.ListAll() store mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
