-- +goose Up
CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    module_id INT,
    contract_id INT NOT NULL,
    created_by UUID NOT NULL,
    assigned_to UUID,
    title VARCHAR(255) NOT NULL,
    message TEXT,
    status VARCHAR(50) DEFAULT 'open',
    gitlab_issue_url TEXT,
    mattermost_thread_url TEXT,
    date_created TIMESTAMP DEFAULT NOW(),
    date_updated TIMESTAMP DEFAULT NOW()
);

DROP TRIGGER IF EXISTS trg_tickets_set_updated ON tickets;
CREATE TRIGGER trg_tickets_set_updated
    BEFORE UPDATE ON tickets
    FOR EACH ROW EXECUTE FUNCTION set_updated_timestamp();

-- +goose Down
DROP TRIGGER IF EXISTS trg_tickets_set_updated ON tickets;
DROP TABLE IF EXISTS tickets;