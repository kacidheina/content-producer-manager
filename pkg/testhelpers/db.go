package testhelpers

import (
	"content-producer-manager/configs"
	"content-producer-manager/pkg/migrations"
	"database/sql"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"log"
	"testing"
)

var connString = fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", configs.DbUser, configs.DbPassword, configs.DbName)

type ID string

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	size     = 32
)

func New() ID {
	return ID(gonanoid.MustGenerate(alphabet, size))
}

func CreateTestDB(t *testing.T, migFile string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Println(nil, "opening DB: %w", err)
		return nil, err
	}

	log.Println(nil, "database connected")

	tmpDB := fmt.Sprintf("tpmdb-%s", New())
	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, tmpDB))
	if err != nil {
		log.Println(nil, "creating temporary DB: %w", err)
		t.Fatal(fmt.Errorf("creating temporary DB: %w", err))
	}

	log.Println(nil, "temporary DB created: %s", tmpDB)

	err = migrations.Up(migFile, connString)
	if err != nil {
		log.Println(nil, "migrating DB: %w", err)
		t.Fatal(fmt.Errorf("migrating DB: %w", err))
	}

	log.Println(nil, "Migration applied")
	return db, nil
}
