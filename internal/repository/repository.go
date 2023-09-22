package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// Prepared stmts совестно с использованием pgxpool не имеет смысла
	// prepareStatementsTimeOut time.Duration = 5 * time.Second

	queryTimeout time.Duration = time.Second
)

type Repository struct {
	conn      *pgxpool.Pool
	initQuery sq.StatementBuilderType
}

func NewRepository(ctx context.Context, conn *pgxpool.Pool) (*Repository, error) {
	r := &Repository{
		conn:      conn,
		initQuery: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}

	// if err := r.prepareStatements(ctx); err != nil {
	// 	return nil, err
	// }

	return r, nil
}

// func (r *Repository) prepareStatements(ctx context.Context) error {

// 	ctx, cf := context.WithTimeout(ctx, prepareStatementsTimeOut)
// 	defer cf()

// r.conn.

// if _, err := r.conn.Prepare(ctx, insertOrderAndDeliveryStmtName, insertOrdersAndDeliverySQL); err != nil {
// 	return err
// }
// if _, err := r.conn.Prepare(ctx, insertPaymentStmtName, insertPaymentSQL); err != nil {
// 	return err
// }

// if _, err := r.conn.Prepare(ctx, selectItemStmtName, selectItemSQL); err != nil {
// 	return err
// }
// if _, err := r.conn.Prepare(ctx, selectOrderPaymentDeliveryStmtName, selectOrderPaymentDeliverySQL); err != nil {
// 	return err
// }

// 	return nil
// }

// func (r *Repository) ReadAllOrdersWithoutItems(ctx context.Context) ([]models.Order, error) {
// 	r.initQuery.Select()
// }
