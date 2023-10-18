package postgres

import (
	"content-producer-manager/pkg/model"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestStoreFile_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	repo := NewRepository(db)

	content := &model.Metadata{
		SenderID:   1,
		ReceiverID: 2,
		FileType:   "txt",
		File:       []byte("Hello, World!"),
		IsPayable:  true,
	}

	// Mock the database query
	mock.ExpectExec(regexp.QuoteMeta(storeFileMetadata)).
		WithArgs(content.SenderID, content.ReceiverID, content.FileType, content.File, content.IsPayable).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.StoreFile(content)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Expectations were not met: %v", err)
	}
}

func TestStoreFile_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	repo := NewRepository(db)

	// Sample metadata
	content := &model.Metadata{
		SenderID:   1,
		ReceiverID: 2,
		FileType:   "txt",
		File:       []byte("Hello, World!"),
		IsPayable:  true,
	}

	mock.ExpectExec(regexp.QuoteMeta(storeFileMetadata)).
		WithArgs(content.SenderID, content.ReceiverID, content.FileType, content.File, content.IsPayable).
		WillReturnError(errors.New("error"))

	err = repo.StoreFile(content)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error")

}
