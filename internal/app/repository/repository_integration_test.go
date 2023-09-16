package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

	"github.com/hablof/order-viewer/internal/database"
	"github.com/hablof/order-viewer/internal/models"
)

const (
	postgresURL string = "postgres://postgres:1234@127.0.0.1:5432/integration_testing?sslmode=disabled"
)

func Test_Repository(t *testing.T) {

	ctx, cf := context.WithTimeout(context.Background(), 15*time.Second)
	cf()

	var c *pgx.Conn
	t.Run("create db connect", func(t *testing.T) {

		var err error
		// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca
		c, err = database.NewPostgres(ctx, postgresURL)
		if err != nil {
			assert.FailNow(t, fmt.Sprintf("failed open db conn: %v", err))
		}
	})
	defer c.Close(ctx)

	t.Run("db setup", func(t *testing.T) {
		if err := dbSetup(ctx, c); err != nil {
			assert.FailNow(t, fmt.Sprintf("failed setup db: %v", err))
		}
	})

	var r *Repository
	t.Run("create repository", func(t *testing.T) {
		var err error
		r, err = NewRepository(context.Background(), c)
		if err != nil {
			assert.FailNow(t, fmt.Sprintf("failed prepare statements: %v", err))
		}
	})

	t.Run("sequential insert", func(t *testing.T) {
		for i, order := range testObjs {
			if err := r.InsertOrder(ctx, order); err != nil {
				assert.FailNow(t, fmt.Sprintf("failed itsert order idx=%d", i))
			}
		}
	})

	t.Run("check count in tables", func(t *testing.T) {
		count := 0
		tables := []string{"orders", "delivery", "payment", "items"}
		expectedCounts := []int{len(testObjs), len(testObjs), len(testObjs), countItems(testObjs)}

		for i, table := range tables {
			if err := c.QueryRow(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
				assert.FailNow(t, "failed select count from orders")
			}
			assert.Equal(t, expectedCounts[i], count, table)
		}

	})
}

func countItems(testObjs []models.Order) int {
	counter := 0
	for _, order := range testObjs {
		counter += len(order.Items)
	}

	return counter
}

func dbSetup(ctx context.Context, c *pgx.Conn) error {

	m, err := migrate.New(
		"file://migrations",
		postgresURL)
	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
