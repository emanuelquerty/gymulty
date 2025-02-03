package postgres

import (
	"context"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/jackc/pgx/v5"
)

var _ domain.TenantStore = (*TenantStore)(nil)

type TenantStore struct {
	conn *pgx.Conn
}

func NewTenantStore(conn *pgx.Conn) *TenantStore {
	return &TenantStore{
		conn: conn,
	}
}

func (t *TenantStore) CreateTenant(ctx context.Context, data domain.Tenant) (domain.Tenant, error) {
	query := "INSERT INTO tenants (business_name, subdomain) VALUES ($1, $2) RETURNING *"

	tx, err := t.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Tenant{}, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, data.BusinessName, data.Subdomain)
	if err != nil {
		return domain.Tenant{}, err
	}

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
