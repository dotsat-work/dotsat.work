-- +goose Up
-- Initial schema for Channel Partner Portal
-- Multi-tenant architecture with tenant isolation

-- ============================================================================
-- TENANTS TABLE
-- Each partner organization is a tenant
-- ============================================================================
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    subdomain TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',  -- active, suspended, inactive
    tier TEXT DEFAULT 'standard',  -- standard, premium, enterprise
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tenants_subdomain ON tenants(subdomain);
CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);

-- ============================================================================
-- USERS TABLE
-- Partner employees - scoped to their tenant
-- ============================================================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    password_hash TEXT NULL,  -- Nullable for passwordless authentication
    role TEXT NOT NULL DEFAULT 'user',  -- admin, user, viewer
    pending_email TEXT NULL,
    email_verified_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, email)  -- Email unique per tenant
);

CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_tenant_email ON users(tenant_id, email);

-- ============================================================================
-- TENANT_SETTINGS TABLE
-- Partner-specific configuration and branding
-- ============================================================================
CREATE TABLE IF NOT EXISTS tenant_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID UNIQUE NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    logo_url TEXT,
    primary_color TEXT DEFAULT '#3B82F6',
    timezone TEXT DEFAULT 'UTC',
    settings JSONB,  -- Flexible key-value storage for custom settings
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tenant_settings_tenant_id ON tenant_settings(tenant_id);

-- ============================================================================
-- TOKENS TABLE
-- Verification tokens, password reset tokens, magic links
-- ============================================================================
CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL,  -- email_verification, password_reset, magic_link
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tokens_token ON tokens(token);
CREATE INDEX IF NOT EXISTS idx_tokens_expires_at ON tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);

-- ============================================================================
-- PROFILES TABLE
-- User profile information (name, bio, avatar, etc.)
-- ============================================================================
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    bio TEXT,
    phone TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id);

-- +goose Down
-- Drop all tables in reverse order (respecting foreign keys)

DROP INDEX IF EXISTS idx_profiles_user_id;
DROP TABLE IF EXISTS profiles;

DROP INDEX IF EXISTS idx_tokens_user_id;
DROP INDEX IF EXISTS idx_tokens_expires_at;
DROP INDEX IF EXISTS idx_tokens_token;
DROP TABLE IF EXISTS tokens;

DROP INDEX IF EXISTS idx_tenant_settings_tenant_id;
DROP TABLE IF EXISTS tenant_settings;

DROP INDEX IF EXISTS idx_users_tenant_email;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_tenant_id;
DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS idx_tenants_status;
DROP INDEX IF EXISTS idx_tenants_subdomain;
DROP TABLE IF EXISTS tenants;
