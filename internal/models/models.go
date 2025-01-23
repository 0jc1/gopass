package models

import (
	"encoding/json"
	"time"
)

type Password struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExportData struct {
	Passwords []Password `json:"passwords"`
	Notes     []Note     `json:"notes"`
}

func (e *ExportData) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *ExportData) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

type SearchResult struct {
	Passwords []Password
	Notes     []Note
}

func NewPassword() *Password {
	now := time.Now()
	return &Password{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewNote() *Note {
	now := time.Now()
	return &Note{
		CreatedAt: now,
		UpdatedAt: now,
	}
}
