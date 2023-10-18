package postgres

import (
	"content-producer-manager/pkg/model"
	"database/sql"
	"errors"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// StoreFile stores the file metadata in the DB
func (r *Repository) StoreFile(content *model.Metadata) error {
	res, err := r.db.Exec(storeFileMetadata,
		content.SenderID,
		content.ReceiverID,
		content.FileType,
		content.File,
		content.IsPayable)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("failed to insert file metadata")
	}

	return nil
}

// Consume gets the content from the DB
func (r *Repository) Consume(senderID int) ([]model.Metadata, error) {
	var binaryData []byte

	rows, err := r.db.Query(getContent, senderID)
	if err != nil {
		return nil, err
	}

	var result []model.Metadata
	for rows.Next() {

		var metadata model.Metadata
		err = rows.Scan(
			&metadata.ID,
			&metadata.SenderID,
			&metadata.ReceiverID,
			&metadata.FileType,
			&binaryData,
			&metadata.IsPayable,
			&metadata.IsPaid,
			&metadata.CreatedAt,
		)

		metadata.File = binaryData

		if err != nil {
			return nil, err
		}

		//simplifying the logic here since the payment service is out of the scope
		//so if the content is payable and not paid, we set it to paid
		if metadata.IsPayable && !metadata.IsPaid {

			err = r.SetAsPaid(metadata.ID)
			if err != nil {
				log.Printf("failed to set as paid, err: %s", err.Error())
				return nil, err
			}
			metadata.IsPaid = true
		}

		result = append(result, metadata)
	}

	return result, nil
}

// SetAsPaid sets the content as paid
func (r *Repository) SetAsPaid(ID int) error {
	res, err := r.db.Exec(setAsPaid, ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("no rows affected")
	}

	return nil

}
