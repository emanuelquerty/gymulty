package postgres

import (
	"context"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ domain.TenantStore = (*TenantStore)(nil)

type TenantStore struct {
	pool *pgxpool.Pool
}

func NewTenantStore(pool *pgxpool.Pool) *TenantStore {
	return &TenantStore{
		pool: pool,
	}
}

func (t *TenantStore) CreateTenant(ctx context.Context, data domain.Tenant) (domain.Tenant, error) {
	query := "INSERT INTO tenants (business_name, subdomain) VALUES ($1, $2) RETURNING *"

	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Tenant{}, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, data.BusinessName, data.Subdomain)
	if err != nil {
		return domain.Tenant{}, err
	}
	defer rows.Close()

	tenant, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Tenant])
	if err != nil {
		return domain.Tenant{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return domain.Tenant{}, err
	}

	return tenant, nil
}
