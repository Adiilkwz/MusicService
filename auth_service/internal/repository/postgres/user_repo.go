package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"auth_service/internal/domain"

	_ "github.com/lib/pq"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) (int64, error) {
	query := `
		INSERT INTO users (email, password_hash, display_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var newID int64
	err := r.db.QueryRowContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		"user",
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}

	return newID, nil
}

func (r *userRepo) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, display_name = $3, avatar_url = $4, role = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`

	var avatarURL sql.NullString
	if user.AvatarURL != "" {
		avatarURL = sql.NullString{String: user.AvatarURL, Valid: true}
	}

	result, err := r.db.ExecContext(ctx, query,
		user.Email, user.PasswordHash, user.DisplayName, avatarURL, user.Role, user.ID,
	)
	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("user not found or no changes made")
	}

	return nil
}

func (r *userRepo) SetResetCode(ctx context.Context, email, code string) error {
	query := `UPDATE users SET reset_code = $1, updated_at = CURRENT_TIMESTAMP WHERE email = $2`

	var resetCode sql.NullString
	if code != "" {
		resetCode = sql.NullString{String: code, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query, resetCode, email)
	if err != nil {
		return fmt.Errorf("repository.SetResetCode: %w", err)
	}
	return nil
}

func (r *userRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, avatar_url, role, reset_code, created_at, updated_at
		FROM users WHERE id = $1
	`
	return r.scanUser(ctx, query, id)
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, avatar_url, role, reset_code, created_at, updated_at
		FROM users WHERE email = $1
	`
	return r.scanUser(ctx, query, email)
}

func (r *userRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, avatar_url, role, reset_code, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository.List: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		var avatarURL, resetCode sql.NullString

		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName,
			&avatarURL, &user.Role, &resetCode, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository.List row scan: %w", err)
		}

		if avatarURL.Valid {
			user.AvatarURL = avatarURL.String
		}
		if resetCode.Valid {
			user.ResetCode = resetCode.String
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepo) scanUser(ctx context.Context, query string, arg any) (*domain.User, error) {
	user := &domain.User{}
	var avatarURL, resetCode sql.NullString

	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName,
		&avatarURL, &user.Role, &resetCode, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("repository.scanUser: %w", err)
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if resetCode.Valid {
		user.ResetCode = resetCode.String
	}

	return user, nil
}
