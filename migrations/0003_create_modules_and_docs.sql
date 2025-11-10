-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS modules (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    responsible_user_id UUID,
    date_created TIMESTAMP DEFAULT NOW(),
    date_updated TIMESTAMP DEFAULT NOW(),
    UNIQUE (project_id, name)
);

CREATE TRIGGER trg_modules_set_updated
    BEFORE UPDATE ON modules
    FOR EACH ROW EXECUTE FUNCTION set_updated_timestamp();

CREATE TABLE IF NOT EXISTS documentations (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    version TEXT,
    uploaded_by UUID,
    date_created TIMESTAMP DEFAULT NOW(),
    date_updated TIMESTAMP DEFAULT NOW(),
    UNIQUE (project_id, file_path)
);

CREATE TRIGGER trg_documentations_set_updated
    BEFORE UPDATE ON documentations
    FOR EACH ROW EXECUTE FUNCTION set_updated_timestamp();
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS documentations CASCADE;
DROP TABLE IF EXISTS modules CASCADE;