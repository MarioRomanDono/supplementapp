-- +goose Up
-- +goose StatementBegin
CREATE TABLE Supplements ( 
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    gtin VARCHAR UNIQUE,
    name VARCHAR,
    brand VARCHAR,
    flavor VARCHAR,
    carbohydrates REAL,
    electrolytes REAL,
    maltodextrose REAL,
    fructose REAL,
    caffeine REAL,
    sodium REAL,
    protein REAL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Supplements;
-- +goose StatementEnd
