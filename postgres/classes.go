package postgres

import (
	"context"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateClass(ctx context.Context, tenantID int, data domain.Class) (domain.Class, error) {
	query :=
		`INSERTO INTO classes (tenant_id, trainer_id, name, description, capacity, starts_at, ends_at)
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
	err = row.Scan(&class.ID, &class.CreatedAt, &class, class.UpdatedAt)
	if err != nil {
		return domain.Class{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return domain.Class{}, err
	}
	return class, nil
}
