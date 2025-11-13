-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_projects (
    user_id UUID NOT NULL,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    contract_id INTEGER REFERENCES contracts(id) ON DELETE SET NULL,
    permissions TEXT,
    assigned_by UUID,
    date_assigned TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, project_id)
);

CREATE INDEX idx_user_projects_user_id ON user_projects(user_id);
CREATE INDEX idx_user_projects_project_id ON user_projects(project_id);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS user_projects CASCADE;