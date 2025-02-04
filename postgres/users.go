package postgres

import (
	"context"
	"errors"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ domain.UserStore = (*UserStore)(nil)

type UserStore struct {
	pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{
		pool: pool,
	}
}

func (u *UserStore) CreateUser(ctx context.Context, tenantID int, data domain.User) (domain.User, error) {
	query := `INSERT INTO users (tenant_id, first_name, last_name, email, password, role) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`

	tx, err := u.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, data.TenantID, data.FirstName, data.LastName, data.Email, data.Password, data.Role)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return domain.User{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserStore) GetUserByID(ctx context.Context, tenantID int, userID int) (domain.User, error) {
	query := "SELECT * FROM users WHERE tenant_id=$1 AND id=$2"
	tx, err := u.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, tenantID, userID)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return domain.User{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserStore) UpdateUser(ctx context.Context, tenantID int, userID int, updates domain.UserUpdate) (domain.User, error) {
	return domain.User{}, errors.New("data store method not implemented")
}

func (u *UserStore) DeleteUserByID(ctx context.Context, tenantID int, userID int) error {
	query := `DELETE FROM users WHERE tenant_id=$1 AND id=$2`
	tx, err := u.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, tenantID, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	return tx.Commit(ctx)
}

func (u *UserStore) GetAllUsers(ctx context.Context, tenantID int) ([]domain.User, error) {
	query := "SELECT * FROM users"
	tx, err := u.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return []domain.User{}, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return []domain.User{}, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return []domain.User{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []domain.User{}, err
	}
	return users, nil
}
