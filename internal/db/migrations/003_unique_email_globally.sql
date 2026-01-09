-- +goose Up
-- Enforce globally unique emails (1 email = 1 tenant)
-- Drop the composite unique constraint and add a global unique constraint

-- Drop the old composite unique constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_tenant_id_email_key;

-- Add global unique constraint on email
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);

-- +goose Down
-- Revert to per-tenant unique emails

-- Drop global unique constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;

-- Restore composite unique constraint
ALTER TABLE users ADD CONSTRAINT users_tenant_id_email_key UNIQUE (tenant_id, email);

