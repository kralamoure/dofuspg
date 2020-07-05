// Package dofuspg implements dofus repository interface for PostgreSQL.
package dofuspg

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var defaultTxOptions = pgx.TxOptions{
	IsoLevel:       pgx.Serializable,
	AccessMode:     pgx.ReadWrite,
	DeferrableMode: pgx.NotDeferrable,
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) (*Repo, error) {
	if pool == nil {
		return nil, errors.New("pool is nil")
	}

	login := &Repo{pool: pool}

	return login, nil
}