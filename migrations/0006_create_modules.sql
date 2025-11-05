-- +goose Up
CREATE TABLE IF NOT EXISTS modules (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    responsible_user_id UUID
);

-- +goose Down
DROP TABLE IF EXISTS modules;
