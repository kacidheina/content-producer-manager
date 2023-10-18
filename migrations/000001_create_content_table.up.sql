CREATE TABLE IF NOT EXISTS content
(
    id serial NOT NULL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    file_type VARCHAR NOT NULL,
    file BYTEA NOT NULL,
    is_payable BOOLEAN NOT NULL,
    is_paid BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);