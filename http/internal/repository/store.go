package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Lafetz/loyalty_marketplace/internal/repository/gen"
	_ "github.com/lib/pq"
)

type Store struct {
	queries *gen.Queries
}

func NewDb(db *sql.DB) *Store {

	queries := gen.New(db)
	return &Store{
		queries: queries,
	}
}
func OpenDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
