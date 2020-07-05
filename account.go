package dofuspg

import (
	"context"
	"time"

	"github.com/kralamoure/dofus"
)

func (r *Repo) CreateAccount(ctx context.Context, account dofus.Account) (id string, err error) {
	query := "INSERT INTO accounts (name, subscription, admin, user_id, last_access, last_ip)" +
		" VALUES ($1, $2, $3, $4, $5, $6)" +
		" RETURNING id;"

	err = repoError(
		r.pool.QueryRow(ctx, query,
			account.Name, account.Subscription, account.Admin, account.UserId, account.LastAccess, account.LastIP).
			Scan(&id),
	)
	return
}

func (r *Repo) Accounts(ctx context.Context) (accounts map[string]dofus.Account, err error) {
	query := "SELECT id, name, subscription, admin, user_id, last_access, last_ip" +
		" FROM accounts;"

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

func (r *Repo) AccountsByUserId(ctx context.Context, userId string) (accounts map[string]dofus.Account, err error) {
	query := "SELECT id, name, subscription, admin, user_id, last_access, last_ip" +
		" FROM accounts" +
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

func (r *Repo) Account(ctx context.Context, id string) (account dofus.Account, err error) {
	query := "SELECT id, name, subscription, admin, user_id, last_access, last_ip" +
		" FROM accounts" +
		" WHERE id = $1;"

	err = repoError(
		r.pool.QueryRow(ctx, query, id).
			Scan(&account.Id, &account.Name, &account.Subscription, &account.Admin, &account.UserId,
				&account.LastAccess, &account.LastIP),
	)
	return
}

func (r *Repo) AccountByName(ctx context.Context, name string) (account dofus.Account, err error) {
	query := "SELECT id, name, subscription, admin, user_id, last_access, last_ip" +
		" FROM accounts" +
		" WHERE name = $1;"

	err = repoError(
		r.pool.QueryRow(ctx, query, name).
			Scan(&account.Id, &account.Name, &account.Subscription, &account.Admin, &account.UserId,
				&account.LastAccess, &account.LastIP),
	)
	return
}

func (r *Repo) SetAccountLastAccessAndLastIP(ctx context.Context, id string, lastAccess time.Time, lastIP string) error {
	query := "UPDATE accounts" +
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
