package ticket_chats

import (
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(chat *TicketChat) error
	GetByID(id int) (*TicketChat, error)
	GetByTicketID(ticketID int) ([]TicketChat, error)
	Update(chat *TicketChat) error
	Delete(id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(chat *TicketChat) error {
	query := `
		INSERT INTO ticket_chats (ticket_id, sender_id, sender_role, message, message_type, mattermost_message_id)
		VALUES (:ticket_id, :sender_id, :sender_role, :message, :message_type, :mattermost_message_id)
		RETURNING id, date_created, date_updated;
	`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(chat, chat)
}

func (r *repository) GetByID(id int) (*TicketChat, error) {
	var chat TicketChat
	err := r.db.Get(&chat, `SELECT * FROM ticket_chats WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *repository) GetByTicketID(ticketID int) ([]TicketChat, error) {
	var chats []TicketChat
	err := r.db.Select(&chats, `SELECT * FROM ticket_chats WHERE ticket_id = $1 ORDER BY date_created ASC`, ticketID)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *repository) Update(chat *TicketChat) error {
	query := `
		UPDATE ticket_chats
		SET message = :message,
		    message_type = :message_type,
		    mattermost_message_id = :mattermost_message_id
		WHERE id = :id
		RETURNING date_updated;
	`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(chat, chat)
}

func (r *repository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM ticket_chats WHERE id = $1`, id)
	return err
}
