package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/vistara-studio/vistara-be/internal/domain/local"
)

var (
	ErrFailedToCommitTransaction   = errors.New("failed to commit transaction")
	ErrFailedToRollbackTransaction = errors.New("failed to rollback transaction")
)

// Repository represents the main repository struct
type Repository struct {
	db *sqlx.DB
}

// RepositoryInterface defines the contract for repository creation
type RepositoryInterface interface {
	NewClient(withTransaction bool) (LocalRepositoryInterface, error)
}

// localRepository implements the local repository with database connection
type localRepository struct {
	queryExecutor namedExtension
}

// LocalRepositoryInterface defines all local business and tourist attraction operations
type LocalRepositoryInterface interface {
	// Transaction management
	Commit() error
	Rollback() error
	
	// Local business operations
	GetAllLocalBusinesses(ctx context.Context, params local.QueryParamRequestGetLocals, out *[]local.Locals) error
	GetLocalBusinessByID(ctx context.Context, data *local.Locals) error
	CreateLocalBusiness(ctx context.Context, business *local.Locals) error
	UpdateLocalBusiness(ctx context.Context, business *local.Locals) error
	DeleteLocalBusiness(ctx context.Context, businessID string) error
	GetReviewsByLocalBusinessID(ctx context.Context, localBusinessID string, out *[]local.Review) error
	
	// Tourist attraction operations
	GetAllTouristAttractions(ctx context.Context, city string, out *[]local.TouristAttractions) error
	GetTouristAttractionByID(ctx context.Context, data *local.TouristAttractions) error
	CreateTouristAttraction(ctx context.Context, attraction *local.TouristAttractions) error
	UpdateTouristAttraction(ctx context.Context, attraction *local.TouristAttractions) error
	DeleteTouristAttraction(ctx context.Context, attractionID string) error
	GetBookingsByTouristAttractionID(ctx context.Context, attractionID string, out *[]local.TourGuideBookings) error
	CreateTourGuideBooking(ctx context.Context, booking *local.TourGuideBookings) error
	GetFullyBookedDates(ctx context.Context, attractionID string, year, month int, dates *[]string) error
}

// namedExtension extends sqlx with named query capabilities
type namedExtension interface {
	sqlx.ExtContext
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
}

// transactionWrapper wraps sqlx.Tx to implement namedExtension
type transactionWrapper struct {
	*sqlx.Tx
}

// NamedExecContext executes a named query with the transaction
func (tw *transactionWrapper) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return sqlx.NamedExecContext(ctx, tw.Tx, query, arg)
}

// NamedQueryContext performs a named query with the transaction
func (tw *transactionWrapper) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return sqlx.NamedQueryContext(ctx, tw.Tx, query, arg)
}

// New creates a new repository instance
func New(database *sqlx.DB) RepositoryInterface {
	return &Repository{db: database}
}

// NewClient creates a new local repository client with optional transaction support
func (r *Repository) NewClient(withTransaction bool) (LocalRepositoryInterface, error) {
	var queryExecutor namedExtension

	queryExecutor = r.db
	if withTransaction {
		transaction, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}
		queryExecutor = &transactionWrapper{transaction}
	}

	return &localRepository{queryExecutor: queryExecutor}, nil
}

// Commit commits the transaction if one exists
func (lr *localRepository) Commit() error {
	switch executor := lr.queryExecutor.(type) {
	case *transactionWrapper:
		return executor.Tx.Commit()
	case *sqlx.DB:
		return nil // No transaction to commit
	default:
		return ErrFailedToCommitTransaction
	}
}

// Rollback rolls back the transaction if one exists
func (lr *localRepository) Rollback() error {
	switch executor := lr.queryExecutor.(type) {
	case *transactionWrapper:
		return executor.Tx.Rollback()
	case *sqlx.DB:
		return nil // No transaction to rollback
	default:
		return ErrFailedToRollbackTransaction
	}
}
