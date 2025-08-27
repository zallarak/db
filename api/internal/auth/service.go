package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/db-xyz/api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type Service struct {
	db        *sql.DB
	jwtSecret []byte
}

func NewService(db *sql.DB) *Service {
	// In production, use a proper secret from environment
	secret := []byte("your-secret-key-change-this")
	return &Service{
		db:        db,
		jwtSecret: secret,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *Service) hashPassword(password string) string {
	salt := []byte("some-salt") // In production, use random salt per password
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return fmt.Sprintf("%x", hash)
}

func (s *Service) Register(email, password string) (*models.User, error) {
	// Check if user exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if count > 0 {
		return nil, ErrUserExists
	}

	// Create user
	user := &models.User{
		ID:        uuid.New().String(),
		Email:     email,
		PwHash:    s.hashPassword(password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO users (id, email, pw_hash, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)`

	_, err = s.db.Exec(query, user.ID, user.Email, user.PwHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *Service) Login(email, password string) (string, *models.User, error) {
	var user models.User
	query := "SELECT id, email, pw_hash, created_at, updated_at FROM users WHERE email = $1"
	
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PwHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return "", nil, ErrInvalidCredentials
	}
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if user.PwHash != s.hashPassword(password) {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, &user, nil
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}