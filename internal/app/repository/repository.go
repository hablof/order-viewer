package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	prepareStatementsTimeOut time.Duration = 5 * time.Second
	queryTimeout             time.Duration = time.Second
)

// Prepared stmts
const (
	// Подготовленные выражения для записи:

	// 1) Для добавления записи (заказ, доставка)
	insertOrderAndDeliveryStmtName string = "insert_order_and_delivery_stmt"
	insertOrdersAndDeliverySQL     string = `
	WITH dq(id) as(
		INSERT INTO delivery (name, phone, zip, city, address, region, email)
		VALUES ($12, $13, $14, $15, $16, $17, $18)
		RETURNING delivery_id
	)
	INSERT INTO orders (
		delivery_id, 
		order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	)
	VALUES 
	(
		(
			SELECT id FROM dq
		), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
	)`

	// 2) Для добавления записи (платёж)
	insertPaymentStmtName string = "insert_payment_stmt"
	insertPaymentSQL      string = `
	INSERT INTO payment 
	(
		transaction, request_id, currency, provider, amount, 
		payment_dt, bank, delivery_cost, goods_total, custom_fee
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	// 3) Для добавления записи в item невозможно подготовить выражение, там переменное число строчек

	// Для получения данных о заказе по его ID:

	// 1) Для получения данных о заказе по его ID
	// selectOrderPaymentDeliveryStmtName = "select_orders_payment_delivery_stmt"
	// selectOrderPaymentDeliverySQL      = `
	// SELECT
	// 	order_uid, track_number, entry, locale, internal_signature,
	// 	customer_id, delivery_service, shardkey, sm_id, date_created,
	// 	oof_shard, name, phone, zip, city, address,	region,	email,
	// 	transaction, request_id, currency, provider, amount, payment_dt,
	// 	bank, delivery_cost, goods_total, custom_fee
	// FROM
	// 	orders
	// 	INNER JOIN delivery USING(delivery_id)
	// 	INNER JOIN payment  ON orders.order_uid = payment.transaction
	// WHERE
	// 	order_uid = $1
	// `
	// //
	// // 2) Для получения списка товаров из заказа по трек-номеру
	// selectItemStmtName = "select_item_stmt"
	// selectItemSQL      = `
	// SELECT
	// 	chrt_id, track_number, price, rid, name, sale,
	// 	size, total_price, nm_id, brand, status
	// FROM
	// 	item
	// WHERE
	// 	track_number = $1
	// `
)

type Repository struct {
	conn      *pgx.Conn
	initQuery sq.StatementBuilderType
}

// Создание репозитория подразумевает формирование
// подготовленных выражений для добавления записи
func NewRepository(ctx context.Context, conn *pgx.Conn) (*Repository, error) {
	r := &Repository{
		conn:      conn,
		initQuery: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}

	if err := r.prepareStatements(ctx); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Repository) prepareStatements(ctx context.Context) error {

	ctx, cf := context.WithTimeout(ctx, prepareStatementsTimeOut)
	defer cf()

	if _, err := r.conn.Prepare(ctx, insertOrderAndDeliveryStmtName, insertOrdersAndDeliverySQL); err != nil {
		return err
	}
	if _, err := r.conn.Prepare(ctx, insertPaymentStmtName, insertPaymentSQL); err != nil {
		return err
	}

	// if _, err := r.conn.Prepare(ctx, selectItemStmtName, selectItemSQL); err != nil {
	// 	return err
	// }
	// if _, err := r.conn.Prepare(ctx, selectOrderPaymentDeliveryStmtName, selectOrderPaymentDeliverySQL); err != nil {
	// 	return err
	// }

	return nil
}

// func (r *Repository) ReadAllOrdersWithoutItems(ctx context.Context) ([]models.Order, error) {
// 	r.initQuery.Select()
// }
