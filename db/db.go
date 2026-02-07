package db

import (
	"context"
	"os"

	"restorapp/db/client"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joho/godotenv"
)

var Queries *client.Queries

func InitDBClient() *pgxpool.Pool {
	godotenv.Load()

	conn, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Error(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Info("Connected to database successfully")

	Queries = client.New(conn)
	log.Info("Created DB client sucessfully")
	return conn
}
