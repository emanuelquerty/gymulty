-- +goose Up
-- +goose StatementBegin
CREATE TABLE classes (
    id SERIAL PRIMARY KEY,
    tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
    trainer_id INT REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR (255),
    description TEXT,
    capacity INT,
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE classes;
-- +goose StatementEnd
