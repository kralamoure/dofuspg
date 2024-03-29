package dofuspg

import (
	"context"
	"time"

	"github.com/kralamoure/dofus"
)

func (r *Db) CreateAccount(ctx context.Context, account dofus.Account) (id string, err error) {
	query := "INSERT INTO dofus.accounts (name, subscription, admin, user_id, last_access, last_ip)" +
		" VALUES ($1, $2, $3, $4, $5, $6)" +
		" RETURNING id;"

	err = dbError(
		r.pool.QueryRow(ctx, query,
			account.Name, account.Subscription, account.Admin, account.UserId, account.LastAccess, account.LastIP).
			Scan(&id),
	)
	return
}

func (r *Db) Accounts(ctx context.Context) (accounts map[string]dofus.Account, err error) {
	query := "SELECT id, name, subscription, admin, user_id, last_access, last_ip" +
		" FROM dofus.accounts;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	accounts = make(map[string]dofus.Account)
	for rows.Next() {
		var account dofus.Account
		err = rows.Scan(&account.Id, &account.Name, &account.Subscription, &account.Admin, &account.UserId,
			&account.LastAccess, &account.LastIP)
		if err != nil {
			return
		}
		accounts[account.Id] = account
	}
	return
}

func (r *Db) AccountsByUserId(ctx context.Context, userId string) (accounts map[string]dofus.Account, err error) {
	query := "SELECT id, name, subscription, admin, user_id, last_access, last_ip" +
		" FROM dofus.accounts" +
		" WHERE user_id = $1;"

	rows, err := r.pool.Query(ctx, query, userId)
	if err != nil {
		return
	}
	defer rows.Close()

	accounts = make(map[string]dofus.Account)
	for rows.Next() {
		var account dofus.Account
		err = rows.Scan(&account.Id, &account.Name, &account.Subscription, &account.Admin, &account.UserId,
			&account.LastAccess, &account.LastIP)
		if err != nil {
			return
		}
		accounts[account.Id] = account
	}
	return
}

func (r *Db) Account(ctx context.Context, id string) (account dofus.Account, err error) {
	query := "SELECT id, name, subscription, admin, user_id, last_access, last_ip" +
		" FROM dofus.accounts" +
		" WHERE id = $1;"

	err = dbError(
		r.pool.QueryRow(ctx, query, id).
			Scan(&account.Id, &account.Name, &account.Subscription, &account.Admin, &account.UserId,
				&account.LastAccess, &account.LastIP),
	)
	return
}

func (r *Db) AccountByName(ctx context.Context, name string) (account dofus.Account, err error) {
	query := "SELECT id, name, subscription, admin, user_id, last_access, last_ip" +
		" FROM dofus.accounts" +
		" WHERE name = $1;"

	err = dbError(
		r.pool.QueryRow(ctx, query, name).
			Scan(&account.Id, &account.Name, &account.Subscription, &account.Admin, &account.UserId,
				&account.LastAccess, &account.LastIP),
	)
	return
}

func (r *Db) SetAccountLastAccessAndLastIP(ctx context.Context, id string, lastAccess time.Time, lastIP string) error {
	query := "UPDATE dofus.accounts" +
		" SET last_access = $2, last_ip = $3" +
		" WHERE id <= $1;"

	tag, err := r.pool.Exec(ctx, query, id, lastAccess, lastIP)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return dofus.ErrNotFound
	}
	return nil
}
