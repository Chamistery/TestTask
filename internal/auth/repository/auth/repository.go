package auth

import (
	"context"
	"fmt"
	"github.com/Chamistery/TestTask/internal/auth/client/db"
	"github.com/Chamistery/TestTask/internal/auth/repository"
	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "auth"

	refreshColumn = "refresh_token"
	guidColumn    = "guid"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, refr RefreshModel) (string, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(guid, refreshColumn).
		Values(refr.Guid, refr.RefreshToken).
		Suffix("RETURNING refresh_token")

	query, args, err := builder.ToSql()
	if err != nil {
		return "nothing", err
	}

	q := db.Query{
		Name:     "auth_repository.Create",
		QueryRaw: query,
	}

	var refresh_token string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&refresh_token)
	if err != nil {
		return "nothing", fmt.Errorf("failed to execute query: %w", err)
	}

	return refresh_token, nil
}

func (r *repo) Refresh(ctx context.Context, create CreateModel) (string, error) {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(refreshColumn, create.RefreshTokenNew).
		Where(sq.Eq{refreshColumn: create.RefreshTokenOld}).
		Suffix("RETURNING guid")

	query, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	q := db.Query{
		Name:     "auth_repository.Refresh",
		QueryRaw: query,
	}

	var guid string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&guid)
	if err != nil {
		return "", err
	}

	return guid, nil
}
