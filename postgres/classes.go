package postgres

import (
	"context"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateClass(ctx context.Context, tenantID int, data domain.Class) (domain.Class, error) {
	query :=
		`INSERT INTO classes (tenant_id, trainer_id, name, description, capacity, starts_at, ends_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Class{}, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, query, data.TenantID, data.TrainerID, data.Name,
		data.Description, data.Capacity, data.StartsAt, data.EndsAt)

	class := data
	err = row.Scan(&class.ID, &class.CreatedAt, &class.UpdatedAt)
	if err != nil {
		return domain.Class{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return domain.Class{}, err
	}
	return class, nil
}

func (s *Store) GetClassByID(ctx context.Context, tenantID int, classID int) (domain.Class, error) {
	query :=
		`SELECT * FROM classes 
		WHERE tenant_id=$1 AND id=$2`
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Class{}, nil
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, tenantID, classID)
	if err != nil {
		return domain.Class{}, nil
	}

	class, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Class])
	if err != nil {
		return domain.Class{}, nil
	}
	err = tx.Commit(ctx)
	if err != nil {
		return domain.Class{}, nil
	}
	return class, nil
}
