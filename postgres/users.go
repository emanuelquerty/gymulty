package postgres

import (
	"context"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateUser(ctx context.Context, tenantID int, data domain.User) (domain.User, error) {
	query :=
		`INSERT INTO users (tenant_id, first_name, last_name, email, password, role) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at`

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, query, data.TenantID, data.FirstName, data.LastName, data.Email, data.Password, data.Role)

	user := data
	err = row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *Store) GetUserByID(ctx context.Context, tenantID int, userID int) (domain.User, error) {
	query := "SELECT * FROM users WHERE tenant_id=$1 AND id=$2"
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, tenantID, userID)
	if err != nil {
		return domain.User{}, err
	}

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

func (s *Store) UpdateUser(ctx context.Context, tenantID int, userID int, updates domain.UserUpdate) (domain.User, error) {
	query, columnValues := buildUserUpdateQuery(tenantID, userID, updates)

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, err
	}
	defer tx.Rollback(ctx)

	rows, err := s.pool.Query(ctx, query, columnValues...)
	if err != nil {
		return domain.User{}, err
	}

	updatedUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return domain.User{}, err
	}
	return updatedUser, nil

}

func (s *Store) DeleteUserByID(ctx context.Context, tenantID int, userID int) error {
	query := `DELETE FROM users WHERE tenant_id=$1 AND id=$2`
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, tenantID, userID)
	if err != nil {
		return err
	}
	rows.Close()

	return tx.Commit(ctx)
}

func (s *Store) GetAllUsers(ctx context.Context, tenantID int) ([]domain.User, error) {
	query := "SELECT * FROM users WHERE tenant_id=$1"
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return []domain.User{}, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, tenantID)
	if err != nil {
		return []domain.User{}, err
	}
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
