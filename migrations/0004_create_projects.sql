-- +goose Up
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    gitlab_project_id INTEGER,
    mattermost_team TEXT
);

-- +goose Down
DROP TABLE IF EXISTS projects;
