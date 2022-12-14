// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: oauth.sql

package db

import (
	"context"
)

const createOauthUser = `-- name: CreateOauthUser :one
INSERT INTO oauths (
    id,
    fullname,
    email,
    provider,
    refresh_token
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, fullname, email, provider, refresh_token, created_at
`

type CreateOauthUserParams struct {
	ID           string `json:"id"`
	Fullname     string `json:"fullname"`
	Email        string `json:"email"`
	Provider     string `json:"provider"`
	RefreshToken string `json:"refresh_token"`
}

func (q *Queries) CreateOauthUser(ctx context.Context, arg CreateOauthUserParams) (Oauths, error) {
	row := q.db.QueryRowContext(ctx, createOauthUser,
		arg.ID,
		arg.Fullname,
		arg.Email,
		arg.Provider,
		arg.RefreshToken,
	)
	var i Oauths
	err := row.Scan(
		&i.ID,
		&i.Fullname,
		&i.Email,
		&i.Provider,
		&i.RefreshToken,
		&i.CreatedAt,
	)
	return i, err
}

const getOauthUser = `-- name: GetOauthUser :one
SELECT id, fullname, email, provider, refresh_token, created_at FROM oauths 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetOauthUser(ctx context.Context, id string) (Oauths, error) {
	row := q.db.QueryRowContext(ctx, getOauthUser, id)
	var i Oauths
	err := row.Scan(
		&i.ID,
		&i.Fullname,
		&i.Email,
		&i.Provider,
		&i.RefreshToken,
		&i.CreatedAt,
	)
	return i, err
}
