package db

import "time"

// User maps 1:1 to the users table.
type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	DisplayName  string    `db:"display_name"`
	Bio          string    `db:"bio"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

// Post maps 1:1 to the posts table.
type Post struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}
