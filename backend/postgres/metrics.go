package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var metrics = struct {
	userCountGauge     prometheus.Gauge
	authCountGauge     prometheus.Gauge
	sessionCountGauge  prometheus.Gauge
	todoItemCountGauge prometheus.Gauge
	todoListCountGauge prometheus.Gauge
}{
	userCountGauge:     promauto.NewGauge(prometheus.GaugeOpts{Name: "todo_db_users", Help: "Total number of users"}),
	authCountGauge:     promauto.NewGauge(prometheus.GaugeOpts{Name: "todo_db_auths", Help: "Total number of auths"}),
	sessionCountGauge:  promauto.NewGauge(prometheus.GaugeOpts{Name: "todo_db_sessions", Help: "Total number of active sessions"}),
	todoItemCountGauge: promauto.NewGauge(prometheus.GaugeOpts{Name: "todo_db_todo_items", Help: "Total number of todo items"}),
	todoListCountGauge: promauto.NewGauge(prometheus.GaugeOpts{Name: "todo_db_todo_lists", Help: "Total number of todo lists"}),
}

func (db *DB) monitorMetrics() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-db.ctx.Done():
			return
		case <-ticker.C:
		}

		if err := db.updateStats(db.ctx); err != nil {
			db.Logger.Errorf("db monitor failed to update stats: %v", err)
		}
	}
}

func (db *DB) updateStats(ctx context.Context) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var n int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&n); err != nil {
		return fmt.Errorf("users count: %v", err)
	}
	metrics.userCountGauge.Set(float64(n))

	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM auths`).Scan(&n); err != nil {
		return fmt.Errorf("auths count: %v", err)
	}
	metrics.authCountGauge.Set(float64(n))

	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions`).Scan(&n); err != nil {
		return fmt.Errorf("sessions count: %v", err)
	}
	metrics.sessionCountGauge.Set(float64(n))

	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM lists`).Scan(&n); err != nil {
		return fmt.Errorf("todo list count: %v", err)
	}
	metrics.todoListCountGauge.Set(float64(n))

	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM items`).Scan(&n); err != nil {
		return fmt.Errorf("todo items count: %v", err)
	}
	metrics.todoItemCountGauge.Set(float64(n))

	return nil
}
