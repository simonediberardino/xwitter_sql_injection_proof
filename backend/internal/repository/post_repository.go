package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"xwitter/internal/db"
	"xwitter/internal/dto"
	"xwitter/internal/mapper"
)

type PostRepository struct {
	pool *pgxpool.Pool
}

func NewPostRepository(pool *pgxpool.Pool) *PostRepository {
	return &PostRepository{pool: pool}
}

func (r *PostRepository) Create(ctx context.Context, post dto.Post) (dto.Post, error) {
	record, err := mapper.PostToDB(post)
	if err != nil {
		return dto.Post{}, err
	}

	err = r.pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, content)
		VALUES ($1, $2)
		RETURNING id, user_id, content, created_at
	`,
		record.UserID,
		record.Content,
	).Scan(
		&record.ID,
		&record.UserID,
		&record.Content,
		&record.CreatedAt,
	)
	if err != nil {
		return dto.Post{}, err
	}

	author, err := scanAuthorByID(ctx, r.pool, record.UserID)
	if err != nil {
		return dto.Post{}, err
	}

	return mapper.PostToDTO(record, author), nil
}

func (r *PostRepository) List(
	ctx context.Context,
	filter dto.PostFilter,
) ([]dto.Post, error) {
	query := `
		SELECT
			posts.id,
			posts.user_id,
			posts.content,
			posts.created_at,
			users.id,
			users.username,
			users.email,
			users.display_name,
			users.bio,
			users.password_hash,
			users.created_at
		FROM posts
		JOIN users ON users.id = posts.user_id
	`

	args := make([]any, 0, 1)

	if filter.UserID != "" {
		query += ` WHERE posts.user_id = $1`

		record, err := mapper.PostToDB(dto.Post{
			UserID: filter.UserID,
		})
		if err != nil {
			return nil, err
		}

		args = append(args, record.UserID)
	}

	query += `
		ORDER BY posts.created_at DESC, posts.id DESC
		LIMIT 100
	`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]db.Post, 0)
	authors := make([]db.User, 0)

	for rows.Next() {
		var post db.Post
		var user db.User

		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&post.CreatedAt,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.DisplayName,
			&user.Bio,
			&user.PasswordHash,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
		authors = append(authors, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mapper.PostsToDTO(posts, authors), nil
}

func scanAuthorByID(
	ctx context.Context,
	pool *pgxpool.Pool,
	userID int,
) (db.User, error) {
	var user db.User

	err := pool.QueryRow(ctx, `
		SELECT
			id,
			username,
			email,
			display_name,
			bio,
			password_hash,
			created_at
		FROM users
		WHERE id = $1
	`,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&user.Bio,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	return user, err
}
