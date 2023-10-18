package migrations

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const maxRetriesOnRefusedConnection = 3
const delayBeforeRetry = 5 * time.Second

func Up(sourceURL, databaseURL string) error {
	var err error
	for try := 0; try < maxRetriesOnRefusedConnection; try++ {
		if err = up(sourceURL, databaseURL); isConnectionRefusedError(err) {
			time.Sleep(delayBeforeRetry)
			continue
		} else {
			return err
		}
	}
	return err
}

func up(sourceURL, databaseURL string) error {
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("unable to setup migration: %w", err)
	}

	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("running migration: %w", err)
	}

	return nil
}

func isConnectionRefusedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "connection refused")
}
