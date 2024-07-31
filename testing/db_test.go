package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var testBun *bun.DB

func TestMain(m *testing.M) {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the test database connection
	testBun = bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(getTestConnStr()))), pgdialect.New())

	// Add verbose logging if needed
	verbose := false
	if verbose {
		testBun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	// Ensure the test table is created before running tests
	if err := createTestTable(); err != nil {
		log.Fatalf("Failed to create test table: %v", err)
	}

	// Run tests
	exitCode := m.Run()

	// Clean up
	if err := dropTestTable(); err != nil {
		log.Fatalf("Failed to drop test table: %v", err)
	}

	// Exit with the test result
	os.Exit(exitCode)
}

func getTestConnStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

func TestDatabaseOperations(t *testing.T) {
	// Insert sample data
	_, err := testBun.ExecContext(context.Background(), `
		INSERT INTO test_users (username, email) VALUES ('john_doe', 'john@example.com')
	`)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Query the data
	var users []User
	err = testBun.NewSelect().Model(&users).Table("test_users").Scan(context.Background())
	if err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}

	// Verify the data
	if len(users) == 0 {
		t.Fatal("Expected to find at least one user, but got none")
	}
	if users[0].Username != "john_doe" {
		t.Fatalf("Expected username to be 'john_doe', but got '%s'", users[0].Username)
	}
}

type User struct {
	ID       int    `bun:"id,pk,autoincrement"`
	Username string `bun:"username,notnull"`
	Email    string `bun:"email,notnull"`
}

func createTestTable() error {
	_, err := testBun.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS test_users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL
		)
	`)
	return err
}

func dropTestTable() error {
	_, err := testBun.ExecContext(context.Background(), `DROP TABLE IF EXISTS test_users`)
	return err
}
