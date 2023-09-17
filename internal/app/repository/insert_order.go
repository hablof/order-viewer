package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/hablof/order-viewer/internal/models"
)

const (
	itemColumns string = "order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status"
)

// Вставка заказа происходит транзакцией в три запроса.
// Запись в orders, delivery, payment блокирующая на уровне таблицы,
// т.к. затрагивает столбцы с ограничением unique. (см. https://postgrespro.ru/docs/postgrespro/15/index-unique-checks)
// Поэтому записать валиндные заказы, но с одинаковыми айди не получится.
func (r *Repository) InsertOrder(ctx context.Context, order models.Order) error {

	// запрос в таблицу item динамический, нужно готовить каждый раз
	query := r.initQuery.Insert("item").Columns(itemColumns)

	for _, item := range order.Items {
		query = query.Values(
			item.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			item.Brand,
			item.Status,
		)
	}

	insertItemQuery, insertItemArgs, err := query.ToSql()
	if err != nil {
		return err
	}

	insertOrderAndDeliveryArgs := r.orderAndDeliveryArgs(order)
	insertPaymentArgs := r.paymentArgs(order.Payment)

	// транзакция не нуждается в повышенном уровне изоляции,
	// т.к. не предполагается конкурентных запросов SELECT
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, insertOrderAndDeliveryStmtName, insertOrderAndDeliveryArgs...); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, insertPaymentStmtName, insertPaymentArgs...); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, insertItemQuery, insertItemArgs...); err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func (r *Repository) paymentArgs(payment models.Payment) []interface{} {
	args := []interface{}{
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDT,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	}

	return args
}

func (r *Repository) orderAndDeliveryArgs(order models.Order) []interface{} {
	args := []interface{}{
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SMID,
		order.DateCreated,
		order.OofShard,

		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	}

	return args
}
