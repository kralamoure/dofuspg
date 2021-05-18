package dofuspg

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/kralamoure/dofus"
)

const (
	errUniqueViolation errCode = "23505"
)

type errCode string

func dbError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", dofus.ErrNotFound, err)
	}

	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return err
	}

	if errCode(pgErr.Code) != errUniqueViolation {
		return err
	}

	var dbErr error
	switch pgErr.ConstraintName {
	case "users_email_key":
		dbErr = dofus.ErrUserEmailAlreadyExists
	case "users_nickname_key":
		dbErr = dofus.ErrUserNicknameAlreadyExists
	case "accounts_name_key":
		dbErr = dofus.ErrAccountNameAlreadyExists
	default:
		dbErr = dofus.ErrAlreadyExists
	}

	return fmt.Errorf("%w: %s", dbErr, err)
}
