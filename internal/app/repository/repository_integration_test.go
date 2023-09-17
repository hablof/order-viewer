package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

	"github.com/hablof/order-viewer/internal/database"
	"github.com/hablof/order-viewer/internal/models"
)

const (
	postgresURL string = "postgres://postgres:1234@127.0.0.1:5432/integration_testing?sslmode=disable"
)

func Test_Repository(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("db setup", func(t *testing.T) {
		if err := dbSetup(t); err != nil {
			assert.FailNow(t, err.Error(), "failed setup db")
		}
	})

	defer dbTeardown(t)

	ctx, cf := context.WithTimeout(context.Background(), 15*time.Second)
	defer cf()

	var c *pgx.Conn
	t.Run("create db connect", func(t *testing.T) {
		var err error
		// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca
		c, err = database.NewPostgres(ctx, postgresURL)
		if err != nil {
			assert.FailNow(t, err.Error(), "failed open db conn")
		}
	})
	defer c.Close(ctx)

	var r *Repository
	t.Run("create repository", func(t *testing.T) {
		var err error
		r, err = NewRepository(context.Background(), c)
		if err != nil {
			assert.FailNow(t, err.Error(), "failed prepare statements")
		}
	})

	t.Run("sequential InsertOrder", func(t *testing.T) {
		for i, order := range testOrders {
			if err := r.InsertOrder(ctx, order); err != nil {
				assert.FailNow(t, err.Error(), fmt.Sprintf("failed itsert order idx=%d", i))
			}
		}
	})

	t.Run("check count in tables", func(t *testing.T) {
		count := 0
		tables := []string{"orders", "delivery", "payment", "item"}
		expectedCounts := []int{len(testOrders), len(testOrders), len(testOrders), countItems(testOrders)}

		for i, table := range tables {
			if err := c.QueryRow(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
				assert.FailNow(t, err.Error(), "failed select count from orders")
			}
			assert.Equal(t, expectedCounts[i], count, table)
		}
	})

	t.Run("ReadAll", func(t *testing.T) {
		mapOrders, err := r.ReadAll(ctx)
		if err != nil {
			assert.FailNow(t, err.Error(), "failed to ReadAll")
		}

		for _, expectedOrder := range testOrders {
			gotOrder, ok := mapOrders[expectedOrder.OrderUID]
			if !ok {
				assert.FailNow(t, "missing order", expectedOrder.OrderUID)
			}

			if err := checkOrdersEqual(t, expectedOrder, gotOrder); err != nil {
				assert.FailNow(t, err.Error(), "comparing order structs failed")
			}
		}
	})
}

func checkOrdersEqual(t *testing.T, exp, got models.Order) error {
	expJSON, err := json.Marshal(exp)
	if err != nil {
		return err
	}
	gotJSON, err := json.Marshal(got)
	if err != nil {
		return err
	}
	if !assert.JSONEq(t, string(expJSON), string(gotJSON)) {
		return errors.New("orders jsons not eq")
	}
	return nil
}

func countItems(testObjs []models.Order) int {
	counter := 0
	for _, order := range testObjs {
		counter += len(order.Items)
	}

	return counter
}

func dbSetup(t *testing.T) error {

	m, err := migrate.New("file://../../../migrations", postgresURL)
	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		return err
	}

	t.Log("db successfully setup")
	return nil
}

func dbTeardown(t *testing.T) error {

	m, err := migrate.New("file://../../../migrations", postgresURL)
	if err != nil {
		return err
	}

	defer m.Close()

	if err := m.Drop(); err != nil {
		return err
	}

	t.Log("db successfully teardown")
	return nil
}
