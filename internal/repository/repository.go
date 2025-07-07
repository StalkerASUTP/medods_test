package repository

import (
	"authentication-app/internal/config"
	"authentication-app/storage/db"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository struct {
	db   *db.Queries
	conn *pgx.Conn
	conf *config.Config
}

func NewRepository(db *db.Queries, conn *pgx.Conn, conf *config.Config) *Repository {
	return &Repository{db: db, conn: conn, conf: conf}
}

func (r *Repository) CreateUser(ctx context.Context, user *db.CreateUserParams) (*db.User, error) {
	const op = "repository.CreateUser"
	trx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer trx.Rollback(ctx)
	trns := r.db.WithTx(trx)
	userNew, err := trns.CreateUser(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("%s: can't create user: %w", op, err)
	}
	if err := trx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: can't commit transaction: %w", op, err)
	}
	return &userNew, nil
}

func (r *Repository) DeactivateUser(ctx context.Context, id pgtype.UUID) error {
	const op = "repository.DeactivateUser"
	trx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer trx.Rollback(ctx)
	trns := r.db.WithTx(trx)
	if err := trns.DeactivateUser(ctx, id); err != nil {
		return fmt.Errorf("%s: can't deactivate user: %w", op, err)
	}
	if err := trx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: can't commit transaction: %w", op, err)
	}
	return nil
}
func (r *Repository) GetUserByID(ctx context.Context, id pgtype.UUID) (*db.User, error) {
	const op = "repository.GetUserByID"
	trx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer trx.Rollback(ctx)
	trns := r.db.WithTx(trx)
	user, err := trns.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: can't get user by id: %w", op, err)
	}
	if err := trx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: can't commit transaction: %w", op, err)
	}
	return &user, nil
}

func (r *Repository) UpdateToken(ctx context.Context, arg db.UpdateTokenParams) (*db.User, error) {
	const op = "repository.UpdateToken"
	trx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer trx.Rollback(ctx)
	trns := r.db.WithTx(trx)
	user, err := trns.UpdateToken(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("%s: can't update token: %w", op, err)
	}
	if err := trx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: can't commit transaction: %w", op, err)
	}
	return &user, nil
}
