package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vistara-studio/vistara-be/internal/domain/session"
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
	NewClient(tx bool) (sessionRepositoryItf, error)
}

type sessionRepository struct {
	q namedExt
}

type sessionRepositoryItf interface {
	Commit() error
	Rollback() error
	CreateSession(ctx context.Context, data session.Table) error
	GetSessionByUserID(ctx context.Context, data *user.Table, out *[]session.Table) error
	DeleteOldestSessionByUserID(ctx context.Context, data session.Table) error
}

type namedExt interface {
	sqlx.ExtContext
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

func New(db *sqlx.DB) RepositoryItf {
	return &repository{db}
}

func (r *repository) NewClient(tx bool) (sessionRepositoryItf, error) {
	var db namedExt

	db = r.DB
	if tx {
		var err error
		db, err = r.DB.Beginx()
		if err != nil {
			return nil, err
		}
	}

	return &sessionRepository{db}, nil
}

func (r *sessionRepository) Commit() error {
	if tx, ok := r.q.(*sqlx.Tx); ok {
		return tx.Commit()
	}

	return errFailedToCommit
}

func (r *sessionRepository) Rollback() error {
	if tx, ok := r.q.(*sqlx.Tx); ok {
		return tx.Rollback()
	}

	return errFailedToRollback
}
