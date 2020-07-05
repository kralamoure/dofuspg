module github.com/kralamoure/dofuspg

go 1.14

require (
	github.com/alexedwards/argon2id v0.0.0-20200522061839-9369edc04b05
	github.com/jackc/pgconn v1.6.1
	github.com/jackc/pgx/v4 v4.7.1
	github.com/kralamoure/dofus v0.0.0-20200705192226-3c9d67956ebf
)

replace github.com/kralamoure/dofus => ../dofus
