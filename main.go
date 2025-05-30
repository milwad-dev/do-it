package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/milwad-dev/do-it/internal/handlers"
	"github.com/milwad-dev/do-it/internal/logger"
	"github.com/milwad-dev/do-it/internal/routers"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
)

var ctx = context.Background()

// @title Do-It Swagger
// @version 1.0
// @description This is the do-it swagger docs
// @termsOfService http://swagger.io/terms/

// @contact.name Do-It Support
// @contact.url https://github.com/milwad-dev
// @contact.email milwad.dev@gmail.com

// @BasePath /api/
func main() {
	// Load environments
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Port
	port := fmt.Sprintf(":%v", os.Getenv("APP_PORT"))

	// Connect to db and set it
	db := connectDB(err)

	defer db.Close()

	// Run tables
	runTables(db)

	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})

	// Ping the Redis server to check the connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	} else {
		fmt.Println("Connected to Redis:", pong)
	}

	// Set and get handler
	handler := handlers.NewDBHandler(db, client)

	// Initialize logger
	isProduction := os.Getenv("APP_ENV") == "production"
	logger.InitLogger(isProduction)
	defer logger.Log.Sync()

	// Router
	r := routers.GetRouter(handler)

	// Serve application
	fmt.Printf("Starting server on %v \n", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}

// connectDB => connect to DB
func connectDB(err error) *sql.DB {
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	dns := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func runTables(db *sql.DB) {
	tables := make(map[string]string)

	tables["users"] = `CREATE TABLE IF NOT EXISTS users (
    	id INT NOT NULL AUTO_INCREMENT UNIQUE,
    	name VARCHAR(255) NOT NULL,
    	email VARCHAR(255) NULL UNIQUE,
    	phone VARCHAR(255) NULL UNIQUE,
    	password VARCHAR(255) NULL,
    	email_verified_at TIMESTAMP NULL DEFAULT NULL,
    	phone_verified_at TIMESTAMP NULL DEFAULT NULL,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	tables["labels"] = `CREATE TABLE IF NOT EXISTS labels (
    	id INT NOT NULL AUTO_INCREMENT UNIQUE,
    	title VARCHAR(255) NOT NULL,
    	color VARCHAR(255) NOT NULL,
    	user_id INT NOT NULL,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`

	tables["tasks"] = `CREATE TABLE IF NOT EXISTS tasks (
		id INT NOT NULL AUTO_INCREMENT UNIQUE,
    	title VARCHAR(255) NOT NULL,
    	description LONGTEXT NOT NULL,
    	status VARCHAR(255) NOT NULL,
    	label_id INT NOT NULL,
    	user_id INT NOT NULL,
    	completed_at TIMESTAMP NULL,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (label_id) REFERENCES labels(id)
    )`

	for _, table := range tables {
		db.Exec(table)
	}
}
