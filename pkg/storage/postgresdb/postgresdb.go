package postgresdb

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func New(constr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			title,
			content,
			author_id,
			created_at,
			published_at
		FROM posts
		JOIN authors ON posts.author_id = authors.id
		ORDER BY id;
	`,
	)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post
	for rows.Next() {
		var t storage.Post
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&t.AuthorID,
			&t.AuthorName,
			&t.CreatedAt,
			&t.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, t)
	}
	return posts, rows.Err()
}

func (s *Store) AddPost(t storage.Post) error {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO posts (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return err
}

func (s *Store) UpdatePost(t storage.Post) error {
	var id int
	err := s.db.QueryRow(context.Background(), `
		UPDATE posts
		SET title = $2 AND content = $3
		WHERE ID = $1;
		`,
		t.ID,
		t.Title,
		t.Content,
	).Scan(&id)
	return err
}

func (s *Store) DeletePost(t storage.Post) error {
	var id int
	err := s.db.QueryRow(context.Background(), `
		DELETE FROM posts
		WHERE ID = $1;
		`,
		t.ID,
	).Scan(&id)
	return err
}
