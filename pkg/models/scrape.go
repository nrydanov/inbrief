package models

import (
	"time"

	"github.com/google/uuid"
)

// BaseRequest аналогично Python BaseRequest с request_id типа UUID или nil (null)
type BaseRequest struct {
	RequestID *uuid.UUID `json:"request_id,omitempty"`
}

// ScrapeRequest наследует BaseRequest (встраивание) и добавляет поля
type ScrapeRequest struct {
	BaseRequest

	ChatFolderLink string    `form:"chat_folder_link" binding:"required,url"`
	RightBound     time.Time `form:"right_bound" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	LeftBound      time.Time `form:"left_bound" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	Social         bool      `form:"social"`
}

// NewScrapeRequest конструктор с установкой значений по умолчанию
func NewScrapeRequest() *ScrapeRequest {
	return &ScrapeRequest{
		ChatFolderLink: "https://t.me/addlist/4eTLcGIrIx9hNzUy",
		RightBound:     time.Now().Local(),
		Social:         false,
	}
}
