package dofuspg

import (
	"context"
	"strings"

	"github.com/kralamoure/dofus"
	"github.com/kralamoure/dofus/dofustyp"
)

func (r *Repo) CreateUser(ctx context.Context, user dofus.User) (id string, err error) {
	query := "INSERT INTO users (email, nickname, gender, community, hash, chat_channels, secret_question, secret_answer)" +
		" VALUES ($1, $2, $3, $4, $5, $6, $7)" +
		" RETURNING id;"

	chatChannels := &strings.Builder{}
	for chatChannel := range user.ChatChannels {
		chatChannels.WriteRune(rune(chatChannel))
	}

	err = repoError(
		r.pool.QueryRow(ctx, query,
			user.Email, user.Nickname, user.Gender, user.Community, user.Hash, chatChannels.String(), user.SecretQuestion, user.SecretAnswer).
			Scan(&id),
	)
	return
}

func (r *Repo) Users(ctx context.Context) (users map[string]dofus.User, err error) {
	query := "SELECT id, email, nickname, gender, community, hash, chat_channels, secret_question, secret_answer" +
		" FROM users;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	users = make(map[string]dofus.User)
	for rows.Next() {
		var user dofus.User
		var chatChannels string

		err = rows.Scan(&user.Id, &user.Email, &user.Nickname, &user.Gender, &user.Community, &user.Hash, &chatChannels,
			&user.SecretQuestion, &user.SecretAnswer)
		if err != nil {
			return
		}

		user.ChatChannels = make(map[dofustyp.ChatChannel]struct{}, len([]rune(chatChannels)))
		for _, chatChannel := range []rune(chatChannels) {
			user.ChatChannels[dofustyp.ChatChannel(chatChannel)] = struct{}{}
		}

		users[user.Id] = user
	}
	return
}

func (r *Repo) User(ctx context.Context, id string) (user dofus.User, err error) {
	query := "SELECT id, email, nickname, gender, community, hash, chat_channels, secret_question, secret_answer" +
		" FROM users" +
		" WHERE id = $1;"

	var chatChannels string

	err = repoError(
		r.pool.QueryRow(ctx, query, id).
			Scan(&user.Id, &user.Email, &user.Nickname, &user.Gender, &user.Community, &user.Hash, &chatChannels, &user.SecretQuestion,
				&user.SecretAnswer),
	)

	user.ChatChannels = make(map[dofustyp.ChatChannel]struct{}, len([]rune(chatChannels)))
	for _, chatChannel := range []rune(chatChannels) {
		user.ChatChannels[dofustyp.ChatChannel(chatChannel)] = struct{}{}
	}

	return
}

func (r *Repo) UserByNickname(ctx context.Context, nickname string) (user dofus.User, err error) {
	query := "SELECT id, email, nickname, gender, community, hash, chat_channels, secret_question, secret_answer" +
		" FROM users" +
		" WHERE nickname = $1;"

	var chatChannels string

	err = repoError(
		r.pool.QueryRow(ctx, query, nickname).
			Scan(&user.Id, &user.Email, &user.Nickname, &user.Gender, &user.Community, &user.Hash, &chatChannels, &user.SecretQuestion,
				&user.SecretAnswer),
	)
	user.ChatChannels = make(map[dofustyp.ChatChannel]struct{}, len([]rune(chatChannels)))
	for _, chatChannel := range []rune(chatChannels) {
		user.ChatChannels[dofustyp.ChatChannel(chatChannel)] = struct{}{}
	}

	return
}

func (r *Repo) UserAddChatChannels(ctx context.Context, id string, chatChannels ...dofustyp.ChatChannel) error {
	tx, err := r.pool.BeginTx(ctx, defaultTxOptions)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var chatChannelsStr string
	err = repoError(tx.QueryRow(ctx, "SELECT chat_channels FROM users WHERE id = $1;", id).
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
		"UPDATE users"+
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
	err = repoError(tx.QueryRow(ctx, "SELECT chat_channels FROM users WHERE id = $1;", id).
		Scan(&chatChannelsStr))
	if err != nil {
		return err
	}

	for _, chatChannel := range chatChannels {
		chatChannelsStr = strings.ReplaceAll(chatChannelsStr, string(chatChannel), "")
	}

	_, err = tx.Exec(ctx,
		"UPDATE users"+
			" SET chat_channels = $2"+
			" WHERE id = $1;", id, chatChannelsStr)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
