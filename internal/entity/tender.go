package entity

import "time"

type (
	Tender struct {
		ID             string    `json:"id" bson:"_id,omitempty"`
		ClientID       string    `json:"client_id" bson:"client_id"`
		Title          string    `json:"title" bson:"title"`
		Description    string    `json:"description" bson:"description"`
		Deadline       time.Time `json:"deadline" bson:"deadline"`
		Budget         float64   `json:"budget" bson:"budget"`
		Status         string    `json:"status" bson:"status"`
		CreatedAt      time.Time `json:"created_at" bson:"created_at"`
		UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
		FileAttachment string    `json:"file_attachment" bson:"file_attachment,omitempty"`
	}
	CreateTenderRequest struct {
		ClientID       string    `json:"client_id"`
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		Deadline       time.Time `json:"deadline"`
		Budget         float64   `json:"budget"`
		FileAttachment string    `json:"file_attachment,omitempty"`
	}
	GetListTender struct {
		Limit int
		Page  int
	}
	UpdateTenderStatusRequest struct {
		ID        string `json:"id"`
		ClientID  string `json:"client_id"`
		NewStatus string `json:"new_status"`
	}
	DeleteTenderRequest struct {
		ID        string `json:"id"`
		ClientID  string `json:"client_id"`
		NewStatus string `json:"new_status"`
	}
)