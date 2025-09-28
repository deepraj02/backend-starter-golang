package store

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Bio          *string   `json:"bio" db:"bio"`
	Provider     string    `json:"provider" db:"provider"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	User         User   `json:"user"`
}

type ForogtPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

type OAuthUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio,omitempty"`
	Provider string `json:"provider"`
}
type PostgresAuthStore struct {
	db *sql.DB
}

func NewPostgresAuthStore(db *sql.DB) *PostgresAuthStore {
	return &PostgresAuthStore{db: db}
}

type AuthStore interface {
	CreateUser(user *User) error
	CreateOAuthUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByEmailAndProvider(email, provider string) (*User, error)
	UpdatePassword(email, hashedPassword string) error
}

func (s *PostgresAuthStore) CreateUser(user *User) error {
	query := `INSERT INTO users (username, email, password_hash, bio)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`
	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresAuthStore) CreateOAuthUser(user *User) error {
	query := `INSERT INTO users (username, email, password_hash, bio, provider)
        VALUES ($1, $2, '', $3, $4)
        RETURNING id, created_at, updated_at`
	err := s.db.QueryRow(query, user.Username, user.Email, user.Bio, user.Provider).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresAuthStore) GetUserByUsername(username string) (*User, error) {
	user := User{}
	query := `
        SELECT id, username, email, password_hash, bio, provider, created_at, updated_at
        FROM users
        WHERE username = $1`

	err := s.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Bio, &user.Provider, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func (s *PostgresAuthStore) GetUserByID(id int64) (*User, error) {
	user := User{}
	query := `
        SELECT id, username, email, password_hash, bio, created_at, updated_at
        FROM users
        WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Bio, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func (s *PostgresAuthStore) GetUserByEmail(email string) (*User, error) {
	user := User{}
	query := `
        SELECT id, username, email, password_hash, bio, provider, created_at, updated_at
        FROM users
        WHERE email = $1`

	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Bio, &user.Provider, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *PostgresAuthStore) GetUserByEmailAndProvider(email, provider string) (*User, error) {
	user := User{}
	query := `
        SELECT id, username, email, password_hash, bio, provider, created_at, updated_at
        FROM users
        WHERE email = $1 AND provider = $2`

	err := s.db.QueryRow(query, email, provider).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Bio, &user.Provider, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *PostgresAuthStore) UpdatePassword(email, hashedPassword string) error {
	query := `
        UPDATE users 
        SET password_hash = $1, updated_at = NOW()
        WHERE email = $2`

	result, err := s.db.Exec(query, hashedPassword, email)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
