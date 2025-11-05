-- +goose Up
CREATE TABLE IF NOT EXISTS documentations (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    version TEXT,
    uploaded_by UUID
);

-- +goose Down
DROP TABLE IF EXISTS documentations;
