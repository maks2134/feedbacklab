package ticket_attachments

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(att *TicketAttachment) error
	GetByID(id int) (*TicketAttachment, error)
	GetByTicketID(ticketID int) ([]TicketAttachment, error)
	Update(att *TicketAttachment) error
	Delete(id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(att *TicketAttachment) error {
	query := `
		INSERT INTO ticket_attachments (ticket_id, file_path, uploaded_by, file_type, description)
		VALUES (:ticket_id, :file_path, :uploaded_by, :file_type, :description)
		RETURNING id, date_created, date_updated;
	`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer func(stmt *sqlx.NamedStmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(stmt)

	return stmt.Get(att, att)
}

func (r *repository) GetByID(id int) (*TicketAttachment, error) {
	var att TicketAttachment
	err := r.db.Get(&att, `SELECT * FROM ticket_attachments WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *repository) GetByTicketID(ticketID int) ([]TicketAttachment, error) {
	var list []TicketAttachment
	err := r.db.Select(&list, `SELECT * FROM ticket_attachments WHERE ticket_id = $1 ORDER BY date_created `, ticketID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) Update(att *TicketAttachment) error {
	query := `
		UPDATE ticket_attachments
		SET file_path = :file_path,
		    file_type = :file_type,
		    description = :description
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
			fmt.Println(err)
		}
	}(stmt)

	return stmt.Get(att, att)
}

func (r *repository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM ticket_attachments WHERE id = $1`, id)
	return err
}
