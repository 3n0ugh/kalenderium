package repository

import (
	"context"
	"database/sql"
	"github.com/3n0ugh/kalenderium/internal/validator"
	"github.com/3n0ugh/kalenderium/pkg/account/database"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("record not found")
)

type User struct {
	UserID       int64  `json:"user_id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash []byte `json:"-"`
}

type AccountRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, email string) (*User, error)
	ServiceStatus(ctx context.Context) error
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(conn database.Connection) AccountRepository {
	return &accountRepository{db: conn.DB()}
}

// Set hashes plain-text password
func (u *User) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

// Matches compare the plain-text password and hashed version
func (u *User) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(plaintextPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// ValidateEmail rules the email validation
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email",
		"must be valid email address")
}

// ValidatePassword rules the password validation
func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

// ValidateUser rules the password and email validations
func ValidateUser(v *validator.Validator, u *User) {
	ValidateEmail(v, u.Email)

	if u.Password != "" {
		ValidatePassword(v, u.Password)
	}

	if u.PasswordHash == nil {
		panic("missing password hash for user")
	}
}

// CreateUser adds given user to mysql database
func (a *accountRepository) CreateUser(ctx context.Context, user User) error {
	query := `INSERT INTO users (email, password) VALUES (?, ?)`

	args := []interface{}{user.Email, user.PasswordHash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := a.db.ExecContext(ctx, query, args)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// GetUser get user from mysql database
func (a *accountRepository) GetUser(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password
        FROM users
        WHERE email = ?
	`
	var user User

	err := a.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

// ServiceStatus a health-check mechanism
func (a *accountRepository) ServiceStatus(ctx context.Context) error {
	return a.db.PingContext(ctx)
}
