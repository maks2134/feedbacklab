package models

import "time"

type UserProject struct {
	ID        int       `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	ProjectID int       `db:"project_id" json:"project_id"`
	Role      string    `db:"role" json:"role"`
	DateAdded time.Time `db:"date_added" json:"date_added"`
}
