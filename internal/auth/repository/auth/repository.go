package auth

import (
	"context"
	"fmt"
	"github.com/Chamistery/TestTask/internal/auth/client/db"
	"github.com/Chamistery/TestTask/internal/auth/model"
	"github.com/Chamistery/TestTask/internal/auth/repository"
	"github.com/Chamistery/TestTask/internal/auth/utils"
	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "auth"

	uuidColumn    = "uuid"
	refreshColumn = "refresh_token"
	guidColumn    = "guid"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, create model.CreateModel) (string, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(uuidColumn, guidColumn, refreshColumn).
		Values(create.Uuid, create.Guid, create.RefreshToken).
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

func (r *repo) Refresh(ctx context.Context, refr model.RefreshModel) (string, error) {
	selectBuilder := sq.
		Select(refreshColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{uuidColumn: refr.Uuid})

	selectQuery, selectArgs, err := selectBuilder.ToSql()
	if err != nil {
		return "", fmt.Errorf("failed to build select query: %w", err)
	}

	selectQ := db.Query{
		Name:     "auth_repository.Select",
		QueryRaw: selectQuery,
	}

	var storedHashedToken string
	err = r.db.DB().QueryRowContext(ctx, selectQ, selectArgs...).Scan(&storedHashedToken)
	if err != nil {
		return "", fmt.Errorf("failed to fetch token record: %w", err)
	}

	if !utils.VerifyHashedToken(storedHashedToken, refr.RefreshTokenOld) {
		return "", fmt.Errorf("Token is wrong")
	}

	updateBuilder := sq.
		Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(refreshColumn, refr.RefreshTokenNew).
		Where(sq.Eq{uuidColumn: refr.Uuid}).
		Suffix("RETURNING " + guidColumn)

	updateQuery, updateArgs, err := updateBuilder.ToSql()
	if err != nil {
		return "", fmt.Errorf("failed to build update query: %w", err)
	}

	updateQ := db.Query{
		Name:     "auth_repository.UpdateRefreshToken",
		QueryRaw: updateQuery,
	}

	var guid string
	err = r.db.DB().QueryRowContext(ctx, updateQ, updateArgs...).Scan(&guid)
	if err != nil {
		return "", fmt.Errorf("failed to update token: %w", err)
	}

	return guid, nil
}
