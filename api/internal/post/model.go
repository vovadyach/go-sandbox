package post

import "time"

type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	USerID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	ImageURL  *string   `json:"image_url"`
}

type WithAuthor struct {
	Post
	AuthorFirstName string `json:"author_first_name"`
	AuthorLastName  string `json:"author_last_name"`
}
