package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     string  `json:"email"`
	Password  string  `json:"-"`
	CreatedAt string  `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := s.db.QueryRowContext(ctx, query, "username", "email", "password").Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `SELECT id, username, first_name, last_name, email, created_at FROM users WHERE id = $1`

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, userID).
		Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil

}
