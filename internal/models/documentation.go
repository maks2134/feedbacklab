package models

import "time"

type Documentation struct {
	ID          int       `db:"id" json:"id"`
	ProjectID   int       `db:"project_id" json:"project_id"`
	Title       string    `db:"title" json:"title"`
	Content     string    `db:"content" json:"content"`
	FileURL     string    `db:"file_url" json:"file_url"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}
