package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgx.Conn
}

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `
    INSERT INTO snippets (title, content, created, expires)
    VALUES($1, $2, CURRENT_TIMESTAMP AT TIME ZONE 'UTC', CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + INTERVAL '1 day' * $3)
    RETURNING id
  `
	id := 0
	err := m.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `
    SELECT id, title, content, created, expires FROM snippets 
    WHERE expires > CURRENT_TIMESTAMP AT TIME ZONE 'UTC' AND id = $1
  `
	s := &Snippet{}
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `
    SELECT id, title, content, created, expires FROM snippets 
    WHERE expires > CURRENT_TIMESTAMP AT TIME ZONE 'UTC' ORDER BY id DESC LIMIT 10
  `
	snippets := []*Snippet{}
	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return snippets, nil
}
