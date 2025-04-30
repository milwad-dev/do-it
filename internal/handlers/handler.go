package handlers

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type DBHandler struct {
	DB          *sql.DB
	redisClient *redis.Client
}

// NewDBHandler => create a new instance from DBHandler and return it
func NewDBHandler(db *sql.DB, redisClient *redis.Client) *DBHandler {
	return &DBHandler{DB: db, redisClient: redisClient}
}
