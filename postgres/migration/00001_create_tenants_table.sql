-- +goose Up
-- +goose StatementBegin
CREATE TABLE tenants (
    id SERIAL PRIMARY KEY,
    business_name VARCHAR (255) UNIQUE NOT NULL,
    subdomain VARCHAR (255) UNIQUE NOT NULL,
    status VARCHAR (50) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tenants;
-- +goose StatementEnd
