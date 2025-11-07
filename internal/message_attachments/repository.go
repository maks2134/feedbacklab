package message_attachments

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(att *MessageAttachment) error
	GetByID(id int) (*MessageAttachment, error)
	GetByChatID(chatID int) ([]MessageAttachment, error)
	Update(att *MessageAttachment) error
	Delete(id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(att *MessageAttachment) error {
	query := `
		INSERT INTO message_attachments (chat_id, file_path, uploaded_by, file_type)
		VALUES (:chat_id, :file_path, :uploaded_by, :file_type)
		RETURNING id, date_created, date_updated;
	`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer func(stmt *sqlx.NamedStmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	return stmt.Get(att, att)
}

func (r *repository) GetByID(id int) (*MessageAttachment, error) {
	var att MessageAttachment
	err := r.db.Get(&att, `SELECT * FROM message_attachments WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *repository) GetByChatID(chatID int) ([]MessageAttachment, error) {
	var list []MessageAttachment
	err := r.db.Select(&list, `SELECT * FROM message_attachments WHERE chat_id = $1 ORDER BY date_created`, chatID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) Update(att *MessageAttachment) error {
	query := `
		UPDATE message_attachments
		SET file_path = :file_path,
		    file_type = :file_type
		WHERE id = :id
		RETURNING date_updated;
	`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer func(stmt *sqlx.NamedStmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	return stmt.Get(att, att)
}

func (r *repository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM message_attachments WHERE id = $1`, id)
	return err
}
