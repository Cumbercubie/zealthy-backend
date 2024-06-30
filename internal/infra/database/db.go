package sql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const SEARCH_PATH = "zealthy"

type DBHandler struct {
	pool *pgxpool.Pool
}

func NewDBHandler(pool *pgxpool.Pool) *DBHandler {
	return &DBHandler{
		pool: pool,
	}
}

func (h *DBHandler) getDbConnection() (*pgxpool.Conn, error) {
	dbConn, err := h.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = dbConn.Exec(context.Background(), fmt.Sprintf("SET search_path TO %s", SEARCH_PATH))
	if err != nil {
		dbConn.Release()
		return nil, err
	}

	return dbConn, nil

}
