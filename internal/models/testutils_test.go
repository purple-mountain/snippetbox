package models

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func newTestDB(t *testing.T) *pgx.Conn {
	db, err := pgx.Connect(context.Background(), "postgres://admin:root@localhost:5432/snippetbox-db")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(context.Background(), string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close(context.Background())
	})

	return db
}
