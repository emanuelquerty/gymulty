-- +goose Up
-- +goose StatementBegin
CREATE TABLE tenants (
    id SERIAL PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    subdomain VARCHAR (255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO tenants (name, subdomain)
VALUES ('FlexBig Gym', 'flexbig.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tenants;
-- +goose StatementEnd
