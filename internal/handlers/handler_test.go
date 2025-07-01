package handlers

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDBHandler(t *testing.T) {
	// Create a mock sql.DB
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a mock redis client (no need to mock connection for this simple test)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Call constructor
	handler := NewDBHandler(db, redisClient)

	// Assertions
	require.NotNil(t, handler)
	require.Equal(t, db, handler.DB)
}
