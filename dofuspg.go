// Package dofuspg is a library that implements dofus.Storer interface (https://github.com/kralamoure/dofus) for a PostgreSQL database.
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

type Db struct {
	pool *pgxpool.Pool
}

func NewDb(pool *pgxpool.Pool) (*Db, error) {
	if pool == nil {
		return nil, errors.New("pool is nil")
	}

	login := &Db{pool: pool}

	return login, nil
}
