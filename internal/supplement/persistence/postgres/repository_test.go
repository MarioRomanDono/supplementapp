package postgres_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/marioromandono/supplementapp/internal/supplement"
	"github.com/marioromandono/supplementapp/internal/supplement/persistence/postgres"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var container *tcpostgres.PostgresContainer
var dbUrl string

const tableName string = "Supplements"

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "supplementapp"
	dbUser := "postgres"
	dbPassword := "password"

	var err error

	container, err = tcpostgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	if err != nil {
		log.Panic(err)
	}

	_, _, err = container.Exec(ctx, []string{
		"psql", "-U", dbUser, "-d", dbName, "-c",
		"CREATE TABLE " + tableName + " ( " +
			"id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, " +
			"gtin VARCHAR UNIQUE, " +
			"name VARCHAR, " +
			"brand VARCHAR, " +
			"flavor VARCHAR, " +
			"carbohydrates REAL, " +
			"electrolytes REAL, " +
			"maltodextrose REAL, " +
			"fructose REAL, " +
			"caffeine REAL, " +
			"sodium REAL, " +
			"protein REAL " +
			")",
	})
	if err != nil {
		log.Panic(err)
	}

	err = container.Snapshot(ctx, tcpostgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		log.Panic(err)
	}

	dbUrl, err = container.ConnectionString(ctx)
	if err != nil {
		log.Panic(err)
	}

	m.Run()
}

func TestPostgresSupplementRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		dbPool := getPool(t, ctx)
		repo := postgres.NewSupplementRepository(dbPool)
		want := supplement.Supplement{
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
		}
		err := repo.Create(ctx, want)

		if err != nil {
			t.Errorf("PostgresSupplementRepository.Create() error = %v, want nil", err)
		}

		got := getSupplement(t, ctx, dbPool, want.Gtin)

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("PostgresSupplementRepository.Create() mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestPostgresSupplementRepository_FindByGtin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		dbPool := getPool(t, ctx)
		repo := postgres.NewSupplementRepository(dbPool)
		got, err := repo.FindByGtin(ctx, "1234567890123")

		if err != nil {
			t.Errorf("PostgresSupplementRepository.FindByGtin() error = %v, want nil", err)
		}

		if got != nil {
			t.Errorf("PostgresSupplementRepository.FindByGtin() got = %v, want nil", got)
		}
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		dbPool := getPool(t, ctx)
		repo := postgres.NewSupplementRepository(dbPool)
		want := supplement.Supplement{
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
		}
		insertSupplement(t, ctx, dbPool, want)

		got, err := repo.FindByGtin(ctx, want.Gtin)

		if err != nil {
			t.Errorf("PostgresSupplementRepository.FindByGtin() error = %v, want nil", err)
		}

		if diff := cmp.Diff(got, &want); diff != "" {
			t.Errorf("PostgresSupplementRepository.FindByGtin() mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestPostgresSupplementRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		dbPool := getPool(t, ctx)
		repo := postgres.NewSupplementRepository(dbPool)
		want := supplement.Supplement{
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
		}
		insertSupplement(t, ctx, dbPool, want)

		want.Name = "new name"
		err := repo.Update(ctx, want)

		if err != nil {
			t.Errorf("PostgresSupplementRepository.Update() error = %v, want nil", err)
		}

		got := getSupplement(t, ctx, dbPool, want.Gtin)

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("PostgresSupplementRepository.Update() mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestPostgresSupplementRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		dbPool := getPool(t, ctx)
		repo := postgres.NewSupplementRepository(dbPool)
		s := supplement.Supplement{
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
		}
		insertSupplement(t, ctx, dbPool, s)

		err := repo.Delete(ctx, s)

		if err != nil {
			t.Errorf("PostgresSupplementRepository.Delete() error = %v, want nil", err)
		}

		got := getSupplement(t, ctx, dbPool, s.Gtin)

		if diff := cmp.Diff(got, supplement.Supplement{}); diff != "" {
			t.Errorf("PostgresSupplementRepository.Delete() mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestPostgresSupplementRepository_ListAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("with no supplements", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		dbPool := getPool(t, ctx)
		repo := postgres.NewSupplementRepository(dbPool)
		got, err := repo.ListAll(ctx)
		want := []supplement.Supplement{}

		if err != nil {
			t.Errorf("PostgresSupplementRepository.ListAll() error = %v, want nil", err)
		}

		if diff := cmp.Diff(got, want, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("PostgresSupplementRepository.ListAll() mismatch (-got +want):\n%s", diff)
		}
	})

	t.Run("with supplements", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(func() {
			err := container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		dbPool := getPool(t, ctx)
		repo := postgres.NewSupplementRepository(dbPool)
		want := []supplement.Supplement{
			{
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
			{
				Gtin:          "1234567890124",
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
		}
		for _, s := range want {
			insertSupplement(t, ctx, dbPool, s)
		}

		got, err := repo.ListAll(ctx)

		if err != nil {
			t.Errorf("PostgresSupplementRepository.ListAll() error = %v, want nil", err)
		}

		if diff := cmp.Diff(got, want, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("PostgresSupplementRepository.ListAll() mismatch (-got +want):\n%s", diff)
		}
	})
}

func getPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()
	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		dbPool.Close()
	})

	return dbPool
}

func insertSupplement(t *testing.T, ctx context.Context, dbPool *pgxpool.Pool, s supplement.Supplement) {
	t.Helper()
	_, err := dbPool.Exec(
		ctx,
		"INSERT INTO "+tableName+
			" (gtin, name, brand, flavor, carbohydrates, electrolytes, maltodextrose, fructose, caffeine, sodium, protein) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		s.Gtin, s.Name, s.Brand, s.Flavor, s.Carbohydrates, s.Electrolytes, s.Maltodextrose, s.Fructose, s.Caffeine, s.Sodium, s.Protein,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func getSupplement(t *testing.T, ctx context.Context, dbPool *pgxpool.Pool, gtin string) supplement.Supplement {
	t.Helper()
	rows, _ := dbPool.Query(
		ctx,
		"SELECT gtin, name, brand, flavor, carbohydrates, electrolytes, maltodextrose, fructose, caffeine, sodium, protein "+
			"FROM "+tableName+" WHERE gtin = $1",
		gtin,
	)
	s, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[supplement.Supplement])

	if err != nil {
		if err == pgx.ErrNoRows {
			return supplement.Supplement{}
		}
		t.Fatal(err)
	}

	return s
}
