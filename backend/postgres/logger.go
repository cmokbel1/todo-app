package postgres

import (
	"context"
	"fmt"

	"github.com/cmokbel1/todo-app/backend/todo"
	"github.com/jackc/pgx/v4"
)

var _ pgx.Logger = (*Logger)(nil)

// Logger implements a pgx.Logger which logs SQL queries
type Logger struct {
	todo.Logger
}

func (l *Logger) Log(ctx context.Context, level pgx.LogLevel, message string, data map[string]interface{}) {
	if message != "Query" {
		return
	}

	var msg string
	if sql, ok := data["sql"].(string); !ok || logBlackListQueries[sql] {
		return
	} else {
		msg = fmt.Sprintf("postgres query %v", sql)
	}

	switch level {
	case pgx.LogLevelDebug:
		l.Debug(msg)
	case pgx.LogLevelInfo:
		l.Info(msg)
	case pgx.LogLevelWarn:
		l.Warn(msg)
	case pgx.LogLevelError:
		l.Error(msg)
	}
}

// logBlackListQueries are queries that are not logged when query logging is enabled.
var logBlackListQueries = map[string]bool{
	// Ignore all metric queries.
	"SELECT COUNT(*) FROM users":                true,
	"SELECT COUNT(*) FROM auths":                true,
	"SELECT COUNT(*) FROM user_app_credentials": true,
	"SELECT COUNT(*) FROM sessions":             true,
	"SELECT COUNT(*) FROM lists":                true,
	"SELECT COUNT(*) FROM items":                true,
}
