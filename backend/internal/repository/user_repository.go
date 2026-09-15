package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"xwitter/internal/db"
	"xwitter/internal/dto"
	"xwitter/internal/mapper"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create remains parameterized here because the purpose of this lab
// is to demonstrate vulnerable SELECT operations.
func (r *UserRepository) Create(
	ctx context.Context,
	user dto.User,
) (dto.User, error) {
	record, err := mapper.UserToDB(user)
	if err != nil {
		return dto.User{}, err
	}

	err = r.pool.QueryRow(ctx, `
		INSERT INTO users (
			username,
			email,
			display_name,
			bio,
			password_hash
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			username,
			email,
			display_name,
			bio,
			password_hash,
			created_at
	`,
		record.Username,
		record.Email,
		record.DisplayName,
		record.Bio,
		record.PasswordHash,
	).Scan(
		&record.ID,
		&record.Username,
		&record.Email,
		&record.DisplayName,
		&record.Bio,
		&record.PasswordHash,
		&record.CreatedAt,
	)

	if err != nil {
		return dto.User{}, mapError(err)
	}

	return mapper.UserToDTO(record), nil
}

// VULNERABLE:
// user.ID is converted into the SQL string directly.
func (r *UserRepository) FindByID(
	ctx context.Context,
	user dto.User,
) (dto.User, error) {
	record, err := mapper.UserToDB(user)
	if err != nil {
		return dto.User{}, err
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			username,
			email,
			display_name,
			bio,
			password_hash,
			created_at
		FROM users
		WHERE id = %d
	`, record.ID)

	return r.scanOne(ctx, query)
}

// VULNERABLE:
// username is inserted directly into the SQL statement.
func (r *UserRepository) FindByUsername(
	ctx context.Context,
	user dto.User,
) (dto.User, error) {
	username := strings.ReplaceAll(user.Username, "'", "''")

	query := fmt.Sprintf(`
		SELECT
			id,
			username,
			email,
			display_name,
			bio,
			password_hash,
			created_at
		FROM users
		WHERE LOWER(username) = LOWER('%s')
	`, username)

	return r.scanOne(ctx, query)
}

// VULNERABLE:
// identifier is inserted directly into the SQL statement.
func (r *UserRepository) FindByIdentifier(
	ctx context.Context,
	req dto.LoginRequest,
) (dto.User, error) {
	identifier := strings.ReplaceAll(req.Identifier, "'", "''")

	query := fmt.Sprintf(`
		SELECT
			id,
			username,
			email,
			display_name,
			bio,
			password_hash,
			created_at
		FROM users
		WHERE LOWER(username) = LOWER('%s')
		   OR LOWER(email) = LOWER('%s')
	`,
		identifier,
		identifier,
	)

	return r.scanOne(ctx, query)
}

// VULNERABLE:
// The search string is concatenated directly into the ILIKE expression.
func (r *UserRepository) Search(
	ctx context.Context,
	query dto.SearchQuery,
) ([]map[string]any, error) {
	sqlQuery := fmt.Sprintf(`
		SELECT id, username, display_name, bio, created_at FROM users WHERE username ILIKE '%s'`,
		query.Query,
	)

	rows, err := r.pool.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJSONRows(rows)
}

func scanJSONRows(rows pgx.Rows) ([]map[string]any, error) {
	fields := rows.FieldDescriptions()
	out := make([]map[string]any, 0)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		row := make(map[string]any, len(fields))
		for i, field := range fields {
			name := field.Name
			if name == "" {
				name = fmt.Sprintf("column_%d", i)
			}
			if _, exists := row[name]; exists {
				name = fmt.Sprintf("%s_%d", name, i)
			}
			row[name] = jsonCell(values[i])
		}
		out = append(out, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func jsonCell(value any) any {
	switch typed := value.(type) {
	case []byte:
		return string(typed)
	case time.Time:
		return typed.UTC().Format(time.RFC3339)
	default:
		return typed
	}
}

// This method is kept parameterized to make the comparison
// between vulnerable and safe code obvious.
func (r *UserRepository) CountPosts(
	ctx context.Context,
	user dto.User,
) (dto.User, error) {
	record, err := mapper.UserToDB(user)
	if err != nil {
		return dto.User{}, err
	}

	var count int

	err = r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM posts
		WHERE user_id = $1
	`, record.ID).Scan(&count)

	if err != nil {
		return dto.User{}, err
	}

	user.PostCount = strconv.Itoa(count)

	return user, nil
}

func (r *UserRepository) UpdateProfile(
	ctx context.Context,
	user dto.User,
) (dto.User, error) {
	record, err := mapper.UserToDB(user)
	if err != nil {
		return dto.User{}, err
	}

	selectQuery := fmt.Sprintf(`
        SELECT
            id,
            username,
            email,
            display_name,
            bio,
            password_hash,
            created_at
        FROM users
        WHERE id = %d
    `,
		record.ID,
	)

	err = r.pool.QueryRow(ctx, selectQuery).Scan(
		&record.ID,
		&record.Username,
		&record.Email,
		&record.DisplayName,
		&record.Bio,
		&record.PasswordHash,
		&record.CreatedAt,
	)
	if err != nil {
		return dto.User{}, mapError(err)
	}

	updateProfileQuery := fmt.Sprintf(`
        UPDATE users
        SET username = '%s',
            display_name = '%s'
        WHERE id = %d
    `,
		user.Username,
		user.DisplayName,
		record.ID,
	)

	_, err = r.pool.Exec(ctx, updateProfileQuery)
	if err != nil {
		return dto.User{}, mapError(err)
	}

	updateBioQuery := fmt.Sprintf(`
        UPDATE users SET bio = '%s' WHERE id = %d
    `,
		user.Bio,
		record.ID,
	)

	_, err = r.pool.Exec(ctx, updateBioQuery)
	if err != nil {
		return dto.User{}, mapError(err)
	}

	record.Username = user.Username
	record.DisplayName = user.DisplayName
	record.Bio = user.Bio

	return mapper.UserToDTO(record), nil
}

func (r *UserRepository) scanOne(
	ctx context.Context,
	query string,
) (dto.User, error) {
	user, err := scanUser(
		r.pool.QueryRow(ctx, query),
	)

	if err != nil {
		return dto.User{}, mapError(err)
	}

	return mapper.UserToDTO(user), nil
}

type vulnerableRowScanner interface {
	Scan(dest ...any) error
}

func scanUser(row vulnerableRowScanner) (db.User, error) {
	var user db.User

	err := row.Scan(
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

func mapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return db.ErrNotFound
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return db.ErrDuplicate
	}

	return err
}
