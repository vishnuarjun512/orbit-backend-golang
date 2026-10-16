package home

import "github.com/jackc/pgx/v5/pgxpool"

type HomeRepo struct {
	db *pgxpool.Pool
}

func NewHomeRepo(db *pgxpool.Pool) *HomeRepo {
	return &HomeRepo{db: db}
}
