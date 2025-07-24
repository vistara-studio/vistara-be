package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vistara-studio/vistara-be/internal/domain/user"
	"github.com/jmoiron/sqlx"
)

var (
	errFailedToCommit   = errors.New("FAILED_TO_COMMIT_TRANSACTION")
	errFailedToRollback = errors.New("FAILED_TO_ROLLBACK_TRANSACTION")
)

type repository struct {
	DB *sqlx.DB
}

type RepositoryItf interface {
	NewClient(tx bool) (userRepositoryItf, error)
}

type userRepository struct {
	q namedExt
}

type userRepositoryItf interface {
	Commit() error
	Rollback() error
	CreateUser(ctx context.Context, data user.Table) error
	GetAccountByEmail(ctx context.Context, data *user.Table) error
}

type namedExt interface {
	sqlx.ExtContext
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

func New(db *sqlx.DB) RepositoryItf {
	return &repository{db}
}

func (r *repository) NewClient(tx bool) (userRepositoryItf, error) {
	var db namedExt

	db = r.DB
	if tx {
		var err error
		db, err = r.DB.Beginx()
		if err != nil {
			return nil, err
		}
	}

	return &userRepository{db}, nil
}

func (r *userRepository) Commit() error {
	if tx, ok := r.q.(*sqlx.Tx); ok {
		return tx.Commit()
	}

	return errFailedToCommit
}

func (r *userRepository) Rollback() error {
	if tx, ok := r.q.(*sqlx.Tx); ok {
		return tx.Rollback()
	}

	return errFailedToRollback
}
