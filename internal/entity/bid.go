package entity

import "time"

type (
	Bid struct {
		ID           string    `json:"id" bson:"_id,omitempty"`
		TenderID     string    `json:"tender_id" bson:"tender_id"`
		ContractorID string    `json:"contractor_id" bson:"contractor_id"`
		Price        float64   `json:"price" bson:"price"`
		DeliveryTime int       `json:"delivery_time" bson:"delivery_time"`
		Comments     string    `json:"comments,omitempty" bson:"comments,omitempty"`
		Status       string    `json:"status" bson:"status"`
		CreatedAt    time.Time `json:"created_at" bson:"created_at"`
		UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	}

	CreateBidRequest struct {
		TenderID     string  `json:"tender_id" bson:"tender_id"`
		ContractorID string  `json:"contractor_id" bson:"contractor_id"`
		Price        float64 `json:"price" bson:"price"`
		DeliveryTime int     `json:"delivery_time" bson:"delivery_time"`
		Comments     string  `json:"comments,omitempty" bson:"comments,omitempty"`
		Status       string  `json:"status" bson:"status"`
	}
	UpdateBidRequest struct {
		TenderID     string  `json:"tender_id" bson:"tender_id"`
		ContractorID string  `json:"contractor_id" bson:"contractor_id"`
		Price        float64 `json:"price" bson:"price"`
		Comments     string  `json:"comments,omitempty" bson:"comments,omitempty"`
		Status       string  `json:"status" bson:"status"`
	}
	GetBidsRequest struct {
		ContractorID string `json:"contractor_id" bson:"contractor_id"`
		Field        string `json:"field" bson:"field"`
		Value        string `json:"value" bson:"value"`
	}
	DeleteBidsRequest struct {
		ContractorID string `json:"contractor_id" bson:"contractor_id"`
	}

	AnnounceWinnerRequest struct {
		ContractorID string `json:"contractor_id" bson:"contractor_id"`
		BidID        string `json:"bid_id" bson:"bid_id"`
	}
)
