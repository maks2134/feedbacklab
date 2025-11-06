-- +goose Up
-- +goose StatementBegin
DO $do$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ticket_status_enum') THEN
CREATE TYPE ticket_status_enum AS ENUM ('open', 'in_progress', 'resolved', 'closed');
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_type_enum') THEN
CREATE TYPE message_type_enum AS ENUM ('text', 'image', 'file', 'system');
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sender_role_enum') THEN
CREATE TYPE sender_role_enum AS ENUM ('client', 'admin');
END IF;
END
$do$;

CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    module_id INT REFERENCES modules(id) ON DELETE SET NULL,
    contract_id INT NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
    created_by UUID NOT NULL,
    assigned_to UUID,
    title VARCHAR(255) NOT NULL,
    message TEXT,
    status ticket_status_enum DEFAULT 'open',
    gitlab_issue_url TEXT,
    mattermost_thread_url TEXT,
    date_created TIMESTAMP DEFAULT NOW(),
    date_updated TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trg_tickets_set_updated
    BEFORE UPDATE ON tickets
    FOR EACH ROW EXECUTE FUNCTION set_updated_timestamp();

CREATE INDEX idx_tickets_project_id ON tickets(project_id);
CREATE INDEX idx_tickets_contract_id ON tickets(contract_id);
CREATE INDEX idx_tickets_created_by ON tickets(created_by);
CREATE INDEX idx_tickets_assigned_to ON tickets(assigned_to);

CREATE TABLE IF NOT EXISTS ticket_chats (
    id SERIAL PRIMARY KEY,
    ticket_id INT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL,
    sender_role sender_role_enum NOT NULL,
    message TEXT NOT NULL,
    message_type message_type_enum DEFAULT 'text',
    mattermost_message_id VARCHAR(100),
    date_created TIMESTAMP DEFAULT NOW(),
    date_updated TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trg_ticket_chats_set_updated
    BEFORE UPDATE ON ticket_chats
    FOR EACH ROW EXECUTE FUNCTION set_updated_timestamp();

CREATE INDEX idx_ticket_chats_ticket_id ON ticket_chats(ticket_id);
CREATE INDEX idx_ticket_chats_sender_id ON ticket_chats(sender_id);

CREATE TABLE IF NOT EXISTS ticket_attachments (
    id SERIAL PRIMARY KEY,
    ticket_id INT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    uploaded_by UUID NOT NULL,
    file_type TEXT,
    description TEXT,
    date_created TIMESTAMP DEFAULT NOW(),
    date_updated TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trg_ticket_attachments_set_updated
    BEFORE UPDATE ON ticket_attachments
    FOR EACH ROW EXECUTE FUNCTION set_updated_timestamp();

CREATE TABLE IF NOT EXISTS message_attachments (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL REFERENCES ticket_chats(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    uploaded_by UUID NOT NULL,
    file_type TEXT,
    date_created TIMESTAMP DEFAULT NOW(),
    date_updated TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trg_message_attachments_set_updated
    BEFORE UPDATE ON message_attachments
    FOR EACH ROW EXECUTE FUNCTION set_updated_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message_attachments CASCADE;
DROP TABLE IF EXISTS ticket_attachments CASCADE;
DROP TABLE IF EXISTS ticket_chats CASCADE;
DROP TABLE IF EXISTS tickets CASCADE;
DROP TYPE IF EXISTS ticket_status_enum;
DROP TYPE IF EXISTS message_type_enum;
DROP TYPE IF EXISTS sender_role_enum;
-- +goose StatementEnd