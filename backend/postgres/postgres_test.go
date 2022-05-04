//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cmokbel1/todo-app/backend/postgres"
)

var (
	host     = env("DB_HOST", "localhost")
	port     = env("DB_PORT", "5432")
	user     = env("DB_USER", "dbuser")
	password = env("DB_PASSWORD", "dbpassword")
	dbname   = env("DB_NAME", "todo")
)

// Ensure the test database can open & close.
func Test_OpenCloseDB(t *testing.T) {
	OpenDB(t)
}

// OpenDB is a utility function that opens a database and creates a separate schema for the specific testing.TB instance.
// It will close the database once the tests have completed.
func OpenDB(tb testing.TB) *postgres.DB {
	tb.Helper()
	rand.Seed(time.Now().UnixNano())

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db := postgres.New(dsn)

	if err := db.Open(context.Background()); err != nil {
		tb.Fatal(err)
	}

	createSchema(tb, db)

	if err := db.Close(); err != nil {
		tb.Fatal(err)
	}

	db = postgres.New(fmt.Sprintf("%s search_path=%s", dsn, schemaName(tb.Name())))

	if err := db.Open(context.Background()); err != nil {
		tb.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		tb.Fatal(err)
	}

	tb.Cleanup(func() {
		if err := dropSchema(tb, db); err != nil {
			tb.Fatal(err)
		}

		if err := db.Close(); err != nil {
			tb.Fatal(err)
		}
	})
	return db
}

// dropSchema drops the schema associated with the current testing.TB
func dropSchema(tb testing.TB, db *postgres.DB) error {
	tx, err := db.BeginTx(context.Background())
	if err != nil {
		tb.Fatal(err)
	}
	defer tx.Rollback()

	schema := schemaName(tb.Name())
	if _, err := tx.ExecContext(context.Background(), fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema)); err != nil {
		tb.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		tb.Fatal(err)
	}
	return nil
}

// createSchema creates a schema for use by the current testing.TB
func createSchema(tb testing.TB, db *postgres.DB) {
	tb.Helper()
	tx, err := db.BeginTx(context.Background())
	if err != nil {
		tb.Fatal(err)
	}
	defer tx.Rollback()

	schema := schemaName(tb.Name())
	if _, err := tx.ExecContext(context.Background(), fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)); err != nil {
		tb.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		tb.Fatal(err)
	}
}

// schemaName creates a safer schema name for use with the DB. It turns names like
// TestUserService_CreateUser/Success into userservice_createuser_success
func schemaName(testName string) string {
	schema := strings.ToLower(testName)
	schema = strings.Replace(schema, "test", "", -1)
	schema = strings.Replace(schema, "/", "_", -1)
	schema = strings.TrimLeft(schema, "_")
	return schema
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// randstr is a utility to generate random strings for various tests reasons, e.g. unique names.
func randstr(n int) *string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	str := string(b)
	return &str
}

func env(key string, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
