package post

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, limit, offset int, status string, sortCol, sortOrder string) ([]WithAuthor, int, error) {
	baseQuery := "FROM posts p JOIN users u ON p.user_id = u.id"
	countQuery := "SELECT COUNT(*) "
	selectQuery := "SELECT p.id, p.created_at, p.updated_at, p.user_id, p.title, p.content, p.status, p.image_url, u.first_name, u.last_name "

	var args []any
	argIndex := 1

	if status != "" {
		baseQuery += fmt.Sprintf(" WHERE p.status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	var total int
	err := r.db.QueryRow(ctx, countQuery+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting posts: %w", err)
	}

	dataQuery := selectQuery + baseQuery +
		fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortCol, sortOrder, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("querying posts: %w", err)
	}
	defer rows.Close()

	var posts []WithAuthor
	for rows.Next() {
		var p WithAuthor
		err := rows.Scan(
			&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.UserID, &p.Title,
			&p.Content, &p.Status, &p.ImageURL, &p.AuthorFirstName, &p.AuthorLastName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning post: %w", err)
		}
		posts = append(posts, p)
	}

	if posts == nil {
		posts = []WithAuthor{}
	}

	return posts, total, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]Post, int, error) {
	var total int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM posts WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting user posts: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, created_at, updated_at, user_id, title, content, status, image_url
		FROM posts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("querying user posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.UserID, &p.Title,
			&p.Content, &p.Status, &p.ImageURL,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning post: %w", err)
		}
		posts = append(posts, p)
	}

	if posts == nil {
		posts = []Post{}
	}

	return posts, total, nil
}
