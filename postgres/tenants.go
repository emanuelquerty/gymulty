package postgres

import (
	"context"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateTenant(ctx context.Context, data domain.Tenant) (domain.Tenant, error) {
	query :=
		`INSERT INTO tenants (business_name, subdomain) 
		VALUES ($1, $2) 
		RETURNING id, status, created_at, updated_at`

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Tenant{}, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, data.BusinessName, data.Subdomain)
	if err != nil {
		return domain.Tenant{}, err
	}
	defer rows.Close()

	tenant, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[domain.Tenant])
	tenant.BusinessName = data.BusinessName
	tenant.Subdomain = data.Subdomain
	if err != nil {
		return domain.Tenant{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return domain.Tenant{}, err
	}

	return tenant, nil
}
