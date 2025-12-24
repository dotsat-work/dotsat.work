package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"dotsat.work/internal/model"
	"dotsat.work/internal/repository"
	"dotsat.work/internal/validation"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailNotVerified   = errors.New("email not verified")
)

type AuthService struct {
	userRepository  repository.UserRepository
	tokenRepository repository.TokenRepository
	jwtSecret       string
	isProduction    bool
	jwtExpiry       time.Duration
}

func NewAuthService(
	userRepository repository.UserRepository,
	tokenRepository repository.TokenRepository,
	jwtSecret string,
	isProduction bool,
	jwtExpiry time.Duration,
) *AuthService {
	return &AuthService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		jwtSecret:       jwtSecret,
		isProduction:    isProduction,
		jwtExpiry:       jwtExpiry,
	}
}

// Login authenticates a user with email and password
func (s *AuthService) Login(email, password string) (*model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	user, err := s.userRepository.ByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("invalid credentials: %w", ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.HasPassword() {
		return nil, fmt.Errorf("this account uses passwordless login. Please use the magic link option")
	}

	err = s.ComparePassword(password, *user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", ErrInvalidCredentials)
	}

	if user.EmailVerifiedAt == nil {
		return nil, fmt.Errorf("email not verified: %w", ErrEmailNotVerified)
	}

	return user, nil
}

// ValidatePassword validates password strength
func (s *AuthService) ValidatePassword(password string) error {
	return validation.ValidatePassword(password)
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a password with a hash
func (s *AuthService) ComparePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateToken generates a random token for magic links, password reset, etc.
func (s *AuthService) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateJWT generates a JWT token for a user
func (s *AuthService) GenerateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID.String(),
		"tenant_id": user.TenantID.String(),
		"email":     user.Email,
		"exp":       time.Now().Add(s.jwtExpiry).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT verifies a JWT token and returns the claims
func (s *AuthService) VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// SetJWTCookie sets the JWT token as an HTTP-only cookie
func (s *AuthService) SetJWTCookie(w http.ResponseWriter, token string, expiry time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  expiry,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearJWTCookie clears the JWT cookie
func (s *AuthService) ClearJWTCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   s.isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}

// SendMagicLink generates a magic link and sends it to the user's email
// For now, this just generates the token - email sending will be added later
func (s *AuthService) SendMagicLink(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	// Validate email
	err := validation.ValidateEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid email: %w", err)
	}

	// Check if a user exists
	user, err := s.userRepository.ByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Delete any existing magic link tokens for this user
	err = s.tokenRepository.DeleteByUserAndType(user.ID, model.TokenTypeMagicLink)
	if err != nil {
		slog.Warn("failed to delete old magic link tokens", "error", err, "user_id", user.ID)
	}

	// Generate magic link token
	magicToken, err := s.GenerateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Save token to the database
	token := &model.Token{
		UserID:    user.ID,
		Type:      model.TokenTypeMagicLink,
		Token:     magicToken,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	err = s.tokenRepository.Create(token)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	slog.Info("magic link generated", "email", user.Email, "token", magicToken)
	return magicToken, nil
}

// VerifyMagicLink verifies the magic link token and returns the authenticated user
func (s *AuthService) VerifyMagicLink(token string) (*model.User, error) {
	// ConsumeToken atomically marks token as used (prevents race conditions)
	tokenModel, err := s.tokenRepository.ConsumeToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired magic link")
	}

	// Verify token type
	if tokenModel.Type != model.TokenTypeMagicLink {
		return nil, fmt.Errorf("invalid token type")
	}

	// Get user
	user, err := s.userRepository.ByID(tokenModel.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Auto-verify email if not already verified
	if user.EmailVerifiedAt == nil {
		now := time.Now()
		user.EmailVerifiedAt = &now
		err = s.userRepository.Update(user)
		if err != nil {
			slog.Warn("failed to verify email", "error", err, "user_id", user.ID)
		}
	}

	slog.Info("user authenticated via magic link", "user_id", user.ID, "email", user.Email)
	return user, nil
}
