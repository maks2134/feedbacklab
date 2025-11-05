-- +goose Up
CREATE TABLE IF NOT EXISTS contracts (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    valid_from DATE,
    valid_to DATE,
    client_company TEXT
);

-- +goose Down
DROP TABLE IF EXISTS contracts;
