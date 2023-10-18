package model

import (
	"errors"
	"time"
)

// Metadata represents the metadata of a content
type Metadata struct {
	ID         int        `json:"id"`
	SenderID   int        `json:"sender_id"`
	ReceiverID int        `json:"receiver_id"`
	File       []byte     `json:"file"`
	FileType   string     `json:"file_type"`
	IsPayable  bool       `json:"is_payable"`
	IsPaid     bool       `json:"is_paid"`
	CreatedAt  *time.Time `json:"created_at"`
}

// Validate check if the metadata is valid
func (c *Metadata) Validate() error {
	if c.SenderID == 0 {
		return errors.New("sender id is missing")
	}
	if len(c.FileType) == 0 {
		return errors.New("file type is missing")
	}
	if c.ReceiverID == 0 {
		return errors.New("receiver id is missing")
	}
	if len(c.File) == 0 {
		return errors.New("file is missing")
	}
	return nil
}
