package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var Bun *bun.DB

const maxRetries = 5
const retryDelay = 2 * time.Second

// Init initializes the PostgreSQL connection using Bun and pgdriver
func Init() {
	var err error

	for i := 0; i < maxRetries; i++ {
		// Construct the database URI in the correct format
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

		// Initialize the SQL database
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))

		// Create a new Bun instance with PostgreSQL dialect
		Bun = bun.NewDB(sqldb, pgdialect.New())

		verbose := false
		if verbose {
			Bun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		}

		// Check the database connection
		err = checkConnection(Bun)
		if err == nil {
			return // Success, exit function
		}

		log.Printf("Failed to connect to database, attempt %d/%d: %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	// If we exhaust all retries, log fatal error
	log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
}

// checkConnection performs a simple query to verify the connection
func checkConnection(db *bun.DB) error {
	var n int
	if err := db.NewSelect().ColumnExpr("1").Scan(context.Background(), &n); err != nil {
		return err
	}
	fmt.Println("Successfully connected to the database")
	return nil
}

// WhereMap adds where clauses to the query builder based on the provided map
func WhereMap(qb bun.QueryBuilder, m fiber.Map) bun.QueryBuilder {
	for k, v := range m {
		qb = qb.Where(fmt.Sprintf("%s = ?", k), v)
	}
	return qb
}
