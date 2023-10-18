package postgres

const (
	storeFileMetadata = `INSERT INTO content (
                     					sender_id,
                   					    receiver_id,
                     					file_type,
                     					file,
                     					is_payable)		
						VALUES ($1, $2, $3, $4, $5)`

	getContent = `	SELECT 	
	    					id,
					sender_id,
					receiver_id,
					file_type,	
					file,
					is_payable,
					is_paid,
					created_at
					FROM content WHERE sender_id = $1`

	setAsPaid = `UPDATE content SET is_paid = true WHERE id = $1`
)
