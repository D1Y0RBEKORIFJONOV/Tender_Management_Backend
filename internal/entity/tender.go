package entity

import "time"

type (
	Tender struct {
		ID             string    `json:"id" bson:"_id,omitempty"`
		ClientID       string    `json:"-" bson:"client_id"`
		Title          string    `json:"title" bson:"title"`
		Description    string    `json:"description" bson:"description"`
		Deadline       time.Time `json:"deadline" bson:"deadline"`
		Budget         float64   `json:"budget" bson:"budget"`
		Status         string    `json:"-" bson:"status"`
		FileAttachment string    `json:"attachment" bson:"file_attachment,omitempty"`
		CreatedAt      time.Time `json:"-" bson:"created_at"`
	}
	CreateTenderRequest struct {
		ClientID       string    `json:"-" bson:"client_id" `
		Title          string    `json:"title" bson:"title" `
		Description    string    `json:"description" bson:"description" `
		Deadline       time.Time `json:"deadline" bson:"deadline" `
		Budget         float64   `json:"budget" bson:"budget"`
		FileAttachment string    `json:"-" bson:"file_attachment,omitempty"`
		CreatedAt      time.Time `json:"-" bson:"created_at"`
	}

	GetListTender struct {
		Limit int
		Page  int
		Field string
		Value string
	}
	UpdateTenderStatusRequest struct {
		ID        string `json:"-"`
		ClientID  string `json:"-"`
		NewStatus string `json:"new_status"`
	}
	DeleteTenderRequest struct {
		ID       string `json:"id"`
		ClientID string `json:"client_id"`
	}
)
