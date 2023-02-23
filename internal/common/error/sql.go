package error

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	DB_ERROR_DUPLICATE_KEY = "23505"
)

func IsDuplicateKeyError(err error) bool {
	if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
		return pgError.Code == DB_ERROR_DUPLICATE_KEY
	}
	return false
}
