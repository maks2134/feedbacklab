-- +goose Up
CREATE TABLE IF NOT EXISTS user_projects (
    user_id UUID NOT NULL,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    contract_id INTEGER REFERENCES contracts(id),
    permissions TEXT,
    uploaded_by UUID,
    PRIMARY KEY (user_id, project_id)
);

-- +goose Down
DROP TABLE IF EXISTS user_projects;
