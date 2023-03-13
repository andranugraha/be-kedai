package error

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	DB_ERROR_DUPLICATE_KEY = "23505"
	DB_ERROR_FOREIGN_KEY   = "23503"
)

func IsDuplicateKeyError(err error) bool {
	if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
		return pgError.Code == DB_ERROR_DUPLICATE_KEY
	}
	return false
}

func IsForeignKeyError(err error) bool {
	if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
		return pgError.Code == DB_ERROR_FOREIGN_KEY
	}
	return false
}
