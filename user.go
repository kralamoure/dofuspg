package dofuspg

import (
	"context"
	"strings"

	"github.com/kralamoure/dofus"
	"github.com/kralamoure/dofus/dofusrepo"
	"github.com/kralamoure/dofus/dofustyp"
)

func (r *Repo) CreateUser(ctx context.Context, user dofus.User) (id string, err error) {
	repoUser, err := dofusrepo.NewUser(user)
	if err != nil {
		return
	}

	query := "INSERT INTO common.users (email, nickname, community, hash, chat_channels, secret_question, secret_answer)" +
		" VALUES ($1, $2, $3, $4, $5, $6, $7)" +
		" RETURNING id;"

	err = repoError(
		r.pool.QueryRow(ctx, query,
			repoUser.Email, repoUser.Nickname, repoUser.Community, repoUser.Hash, repoUser.ChatChannels, repoUser.SecretQuestion, repoUser.SecretAnswer).
			Scan(&id),
	)
	return
}

func (r *Repo) Users(ctx context.Context) (users map[string]dofus.User, err error) {
	query := "SELECT id, email, nickname, community, hash, chat_channels, secret_question, secret_answer" +
		" FROM common.users;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	repoUsers := make(map[string]dofusrepo.User)
	for rows.Next() {
		var repoUser dofusrepo.User
		err = rows.Scan(&repoUser.Id, &repoUser.Email, &repoUser.Nickname, &repoUser.Community, &repoUser.Hash, &repoUser.ChatChannels,
			&repoUser.SecretQuestion, &repoUser.SecretAnswer)
		if err != nil {
			return
		}
		repoUsers[repoUser.Id] = repoUser
	}

	users = make(map[string]dofus.User, len(repoUsers))
	for k := range repoUsers {
		user, err2 := repoUsers[k].Entity()
		if err2 != nil {
			err = err2
			return
		}
		users[k] = user
	}
	return
}

func (r *Repo) User(ctx context.Context, id string) (user dofus.User, err error) {
	var repoUser dofusrepo.User

	query := "SELECT id, email, nickname, community, hash, chat_channels, secret_question, secret_answer" +
		" FROM common.users" +
		" WHERE id = $1;"

	err = repoError(
		r.pool.QueryRow(ctx, query, id).
			Scan(&repoUser.Id, &repoUser.Email, &repoUser.Nickname, &repoUser.Community, &repoUser.Hash, &repoUser.ChatChannels, &repoUser.SecretQuestion,
				&repoUser.SecretAnswer),
	)
	return repoUser.Entity()
}

func (r *Repo) UserByNickname(ctx context.Context, nickname string) (user dofus.User, err error) {
	var repoUser dofusrepo.User

	query := "SELECT id, email, nickname, community, hash, chat_channels, secret_question, secret_answer" +
		" FROM common.users" +
		" WHERE nickname = $1;"

	err = repoError(
		r.pool.QueryRow(ctx, query, nickname).
			Scan(&repoUser.Id, &repoUser.Email, &repoUser.Nickname, &repoUser.Community, &repoUser.Hash, &repoUser.ChatChannels, &repoUser.SecretQuestion,
				&repoUser.SecretAnswer),
	)
	return repoUser.Entity()
}

func (r *Repo) UserAddChatChannels(ctx context.Context, id string, chatChannels ...dofustyp.ChatChannel) error {
	tx, err := r.pool.BeginTx(ctx, defaultTxOptions)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var chatChannelsStr string
	err = repoError(tx.QueryRow(ctx, "SELECT chat_channels FROM common.users WHERE id = $1;", id).
		Scan(&chatChannelsStr))
	if err != nil {
		return err
	}

	sb := &strings.Builder{}
	sb.WriteString(chatChannelsStr)
	for _, chatChannel := range chatChannels {
		if !strings.ContainsRune(chatChannelsStr, rune(chatChannel)) {
			sb.WriteRune(rune(chatChannel))
		}
	}

	_, err = tx.Exec(ctx,
		"UPDATE common.users"+
			" SET chat_channels = $2"+
			" WHERE id = $1;", id, sb.String())
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repo) UserRemoveChatChannels(ctx context.Context, id string, chatChannels ...dofustyp.ChatChannel) error {
	tx, err := r.pool.BeginTx(ctx, defaultTxOptions)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var chatChannelsStr string
	err = repoError(tx.QueryRow(ctx, "SELECT chat_channels FROM common.users WHERE id = $1;", id).
		Scan(&chatChannelsStr))
	if err != nil {
		return err
	}

	for _, chatChannel := range chatChannels {
		chatChannelsStr = strings.ReplaceAll(chatChannelsStr, string(chatChannel), "")
	}

	_, err = tx.Exec(ctx,
		"UPDATE common.users"+
			" SET chat_channels = $2"+
			" WHERE id = $1;", id, chatChannelsStr)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
