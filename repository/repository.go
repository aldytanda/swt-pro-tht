// This file contains the repository implementation layer.
package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository struct {
	Db           *sql.DB
	JWTSecretKey string
}

type NewRepositoryOptions struct {
	Dsn          string
	JWTSecretKey string
}

func NewRepository(opts NewRepositoryOptions) *Repository {
	db, err := sql.Open("postgres", opts.Dsn)
	if err != nil {
		panic(err)
	}
	return &Repository{
		Db:           db,
		JWTSecretKey: opts.JWTSecretKey,
	}
}
