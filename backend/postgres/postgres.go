package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/cmokbel1/todo-app/backend/todo"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DB is a database driver wrapper which exposes utility methods for interacting with Postgres.
type DB struct {
	db     *sqlx.DB
	ctx    context.Context
	cancel func()

	// Connection string
	DSN string
	// Application logger
	Logger todo.Logger
	// EnableQueryLogging toggles INFO logging of underlying SQL queries.
	EnableQueryLogging bool

	// Now returns current time in UTC rounded to the nearest microsecond
	Now func() time.Time
}

func New(dsn string) *DB {
	db := &DB{
		DSN:    dsn,
		Now:    func() time.Time { return time.Now().UTC().Round(time.Microsecond) },
		Logger: todo.NewLogger(),
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open(ctx context.Context) error {
	if db.DSN == "" {
		return errors.New("dsn is required")
	}

	cfg, err := pgx.ParseConfig(db.DSN)
	if err != nil {
		return err
	}
	if db.EnableQueryLogging {
		cfg.Logger = &Logger{db.Logger}
	}
	sqlDB := stdlib.OpenDB(*cfg)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 3)

	if db.db = sqlx.NewDb(sqlDB, "pgx"); err != nil {
		return fmt.Errorf("failed to open database driver: %v", err)
	}

	db.Logger.Info("pinging database")

	if err = db.ping(ctx); err != nil {
		return err
	}

	db.Logger.Info("successfully pinged database")
	return nil
}

func (db *DB) Migrate() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	goose.SetBaseFS(migrationsFS)
	goose.SetLogger(db.Logger)

	if err := goose.Up(db.db.DB, "migrations"); err != nil {
		return err
	}

	go db.monitorMetrics()
	return nil
}

func (db *DB) Close() error {
	db.cancel()

	if db.db != nil {
		return db.db.Close()
	}

	return nil
}

func (db *DB) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now(),
	}, nil
}

func (db *DB) ping(ctx context.Context) error {
	pctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	for {
		if err := db.db.PingContext(pctx); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return errors.New("failed to ping database within max allotted ping time")
			}
			db.Logger.Warnf("failed to ping database: %v", err)
		} else {
			return nil
		}

		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					return errors.New("database ping interrupted")
				}
				return err
			}
		case <-time.After(time.Second * 3):
		}
	}
}

// Tx is a transaction wrapper with configurable now time parameter.
type Tx struct {
	*sql.Tx
	db  *DB
	now time.Time
}

// Time is a helper type used on time.Time to ensure that records read/written to postgres are
// properly formatted and in UTC time rounded to the nearest microsecond.
type Time time.Time

func (t *Time) Value() (driver.Value, error) {
	if t == nil || (*time.Time)(t).IsZero() {
		return nil, nil
	}
	return (*time.Time)(t).UTC().Round(time.Microsecond), nil
}

// Scan reads a time value from the database.
func (t *Time) Scan(value interface{}) error {
	if value == nil {
		*(*time.Time)(t) = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case *time.Time:
		*(*time.Time)(t) = v.UTC().Round(time.Microsecond)
		return nil
	case time.Time:
		*(*time.Time)(t) = v.UTC().Round(time.Microsecond)
		return nil
	}
	return fmt.Errorf("postgres/Time.Scan: cannot scan %T to time.Time", value)
}

// FormatLimitOffset returns a LIMIT/OFFSET clause or an empty string if none
// is specified.
func FormatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	} else if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}
	return ""
}

func NewSessionStore(db *DB) *postgresstore.PostgresStore {
	return postgresstore.NewWithCleanupInterval(db.db.DB, time.Minute*30)
}
