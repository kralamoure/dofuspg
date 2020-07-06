module github.com/kralamoure/dofuspg

go 1.14

require (
	github.com/jackc/pgconn v1.6.1
	github.com/jackc/pgx/v4 v4.7.1
	github.com/kralamoure/dofus v0.0.0-20200705225418-6b7bc89b411c
)

replace github.com/kralamoure/dofus => ../dofus
