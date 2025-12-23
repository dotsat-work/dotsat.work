-- ============================================================================
-- TOKENS TABLE
-- For email verification, password reset, magic links, email change
-- ============================================================================
CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL,  -- 'email_verify', 'password_reset', 'magic_link', 'email_change'
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for fast token lookup
CREATE INDEX IF NOT EXISTS idx_tokens_token ON tokens(token);

-- Index for user lookup
CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);

-- Index for cleanup queries (find expired/used tokens)
CREATE INDEX IF NOT EXISTS idx_tokens_expires_at ON tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_tokens_used_at ON tokens(used_at);

