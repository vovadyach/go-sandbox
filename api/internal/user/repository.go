package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, limit, offset int, role, status string, sortCol, sortOrder string) ([]User, int, error) {
	baseQuery := "FROM users"
	countQuery := "SELECT COUNT(*) "
	selectQuery := "SELECT id, created_at, updated_at, first_name, last_name, email, role, status, country, avatar_url "

	var args []any
	var conditions []string
	argIndex := 1

	if role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, role)
		argIndex++
	}

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			baseQuery += " AND " + conditions[i]
		}
	}

	// Count
	var total int
	err := r.db.QueryRow(ctx, countQuery+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting: %w", err)
	}

	// Query
	dataQuery := selectQuery + baseQuery +
		fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortCol, sortOrder, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("querying users: %w", err)
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FirstName,
			&user.LastName, &user.Email, &user.Role, &user.Status, &user.Country, &user.AvatarURL,
		); err != nil {
			return nil, 0, fmt.Errorf("scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating users: %w", err)
	}

	if users == nil {
		users = []User{}
	}

	return users, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, created_at, updated_at, first_name, last_name, email, role, status, country, avatar_url 
		FROM users 
		WHERE id = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FirstName,
		&user.LastName, &user.Email, &user.Role, &user.Status, &user.Country, &user.AvatarURL,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying: %s: %w", id, err)
	}

	return &user, nil
}
