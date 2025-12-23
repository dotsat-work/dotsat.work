package ctxkeys

import (
	"context"

	"dotsat.work/internal/config"
	"dotsat.work/internal/model"
)

// contextKey is a type for context keys to avoid collisions
type contextKey string

const (
	UserKey      contextKey = "user"
	TenantKey    contextKey = "tenant"
	ConfigKey    contextKey = "config"
	CSRFTokenKey contextKey = "csrf_token"
)

// User retrieves the user from context
func User(ctx context.Context) *model.User {
	user, _ := ctx.Value(UserKey).(*model.User)
	return user
}

// WithUser adds a user to the context
func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

// Tenant retrieves the tenant from context
func Tenant(ctx context.Context) *model.Tenant {
	tenant, _ := ctx.Value(TenantKey).(*model.Tenant)
	return tenant
}

// WithTenant adds a tenant to the context
func WithTenant(ctx context.Context, tenant *model.Tenant) context.Context {
	return context.WithValue(ctx, TenantKey, tenant)
}

// Config retrieves the config from context
func Config(ctx context.Context) *config.Config {
	cfg, _ := ctx.Value(ConfigKey).(*config.Config)
	return cfg
}

// WithConfig adds config to the context
func WithConfig(ctx context.Context, cfg *config.Config) context.Context {
	return context.WithValue(ctx, ConfigKey, cfg)
}

// CSRFToken retrieves the CSRF token from context
func CSRFToken(ctx context.Context) string {
	token, _ := ctx.Value(CSRFTokenKey).(string)
	return token
}

// WithCSRFToken adds a CSRF token to the context
func WithCSRFToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, CSRFTokenKey, token)
}
