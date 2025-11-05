-- +goose Up
CREATE TABLE IF NOT EXISTS ticket_chats (
    id SERIAL PRIMARY KEY,
    ticket_id INT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL,
    sender_role VARCHAR(50) NOT NULL CHECK (sender_role IN ('client','support')),
    message TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text',
    mattermost_message_id VARCHAR(100),
    date_created TIMESTAMP DEFAULT NOW(),
    date_updated TIMESTAMP DEFAULT NOW()
    );

DROP TRIGGER IF EXISTS trg_ticket_chats_set_updated ON ticket_chats;
CREATE TRIGGER trg_ticket_chats_set_updated
    BEFORE UPDATE ON ticket_chats
    FOR EACH ROW EXECUTE FUNCTION set_updated_timestamp();

CREATE TABLE IF NOT EXISTS ticket_attachments (
    id SERIAL PRIMARY KEY,
    ticket_id INT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    uploaded_by UUID NOT NULL,
    file_type VARCHAR(50),
    description TEXT,
    date_created TIMESTAMP DEFAULT NOW(),
    date_updated TIMESTAMP DEFAULT NOW()
    );

DROP TRIGGER IF EXISTS trg_ticket_attachments_set_updated ON ticket_attachments;
CREATE TRIGGER trg_ticket_attachments_set_updated
    BEFORE UPDATE ON ticket_attachments
    FOR EACH ROW EXECUTE FUNCTION set_updated_timestamp();

CREATE TABLE IF NOT EXISTS message_attachments (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL REFERENCES ticket_chats(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    uploaded_by UUID NOT NULL,
    file_type VARCHAR(50),
    date_created TIMESTAMP DEFAULT NOW(),
    date_updated TIMESTAMP DEFAULT NOW()
    );

DROP TRIGGER IF EXISTS trg_message_attachments_set_updated ON message_attachments;
CREATE TRIGGER trg_message_attachments_set_updated
    BEFORE UPDATE ON message_attachments
    FOR EACH ROW EXECUTE FUNCTION set_updated_timestamp();

-- +goose Down
DROP TABLE IF EXISTS message_attachments;
DROP TABLE IF EXISTS ticket_attachments;
DROP TABLE IF EXISTS ticket_chats;