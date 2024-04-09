package postgres

import (
	"context"

	"github.com/marioromandono/supplementapp/internal/supplement"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSupplementRepository struct {
	db        *pgxpool.Pool
	tableName string
}

func NewSupplementRepository(db *pgxpool.Pool) *PostgresSupplementRepository {
	return &PostgresSupplementRepository{db: db, tableName: "Supplements"}
}

func (r *PostgresSupplementRepository) FindByGtin(ctx context.Context, gtin string) (*supplement.Supplement, error) {
	rows, _ := r.db.Query(
		ctx,
		"SELECT gtin, name, brand, flavor, carbohydrates, electrolytes, maltodextrose, fructose, caffeine, sodium, protein "+
			"FROM "+r.tableName+" WHERE gtin = $1",
		gtin,
	)
	s, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[supplement.Supplement])

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return s, err
}

func (r *PostgresSupplementRepository) Create(ctx context.Context, s supplement.Supplement) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO "+r.tableName+
			" (gtin, name, brand, flavor, carbohydrates, electrolytes, maltodextrose, fructose, caffeine, sodium, protein) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		s.Gtin, s.Name, s.Brand, s.Flavor, s.Carbohydrates, s.Electrolytes, s.Maltodextrose, s.Fructose, s.Caffeine, s.Sodium, s.Protein,
	)

	return err
}

func (r *PostgresSupplementRepository) Update(ctx context.Context, s supplement.Supplement) error {
	_, err := r.db.Exec(
		ctx,
		"UPDATE "+r.tableName+
			" SET name = $1, brand = $2, flavor = $3, carbohydrates = $4, electrolytes = $5, maltodextrose = $6, fructose = $7, caffeine = $8, sodium = $9, protein = $10 "+
			"WHERE gtin = $11",
		s.Name, s.Brand, s.Flavor, s.Carbohydrates, s.Electrolytes, s.Maltodextrose, s.Fructose, s.Caffeine, s.Sodium, s.Protein, s.Gtin,
	)

	return err
}

func (r *PostgresSupplementRepository) Delete(ctx context.Context, s supplement.Supplement) error {
	_, err := r.db.Exec(ctx, "DELETE FROM "+r.tableName+" WHERE gtin = $1", s.Gtin)
	return err
}

func (r *PostgresSupplementRepository) ListAll(ctx context.Context) ([]supplement.Supplement, error) {
	rows, _ := r.db.Query(
		ctx,
		"SELECT gtin, name, brand, flavor, carbohydrates, electrolytes, maltodextrose, fructose, caffeine, sodium, protein "+
			"FROM "+r.tableName,
	)
	return pgx.CollectRows(rows, pgx.RowToStructByName[supplement.Supplement])
}
