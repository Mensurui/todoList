package data

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type UserModel struct {
	DB *sql.DB
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)

	if err != nil {
		return err
	}

	p.plaintext = &plainTextPassword
	p.hash = hash

	return nil
}

func (p password) Match(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users(name, email, password_hash, activated)
	VALUES($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	args := []interface{}{&user.Name, &user.Email, &user.Password.hash, &user.Activated}

	err := u.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}

	}
	return nil
}

func (u *UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, name, email, password_hash, activated
	FROM users
	WHERE email=$1`

	var user User

	err := u.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password.hash, &user.Activated)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModel) GetByToken(tokenPlainText string) (*User, error) {
	// Hash the incoming token
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	// Log the hashed token for debugging
	log.Printf("Token Hash: %x", tokenHash)

	query := `
    SELECT users.id, users.name, users.email, users.password_hash, users.activated
    FROM users
    INNER JOIN tokens
    ON users.id = tokens.user_id
    WHERE tokens.hash = $1
    AND tokens.expiry > $2
    `

	var user User
	// Convert the token hash to a string representation
	args := []interface{}{tokenHash[:], time.Now()}

	// Execute the query
	err := u.DB.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.hash, // Make sure this matches your User struct
		&user.Activated,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
