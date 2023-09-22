package repository

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hablof/order-viewer/internal/models"
	"github.com/hablof/order-viewer/internal/service"
)

const (
	insertOrdersAndDeliverySQL string = `
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

	insertPaymentSQL string = `
	INSERT INTO payment 
	(
		transaction, request_id, currency, provider, amount, 
		payment_dt, bank, delivery_cost, goods_total, custom_fee
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
)

const (
	itemColumns string = "order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status"
)

// Вставка заказа происходит транзакцией в три запроса.
// Запись в orders, delivery, payment блокирующая на уровне таблицы,
// т.к. затрагивает столбцы с ограничением unique. (см. https://postgrespro.ru/docs/postgrespro/15/index-unique-checks)
// Поэтому записать валиндные заказы, но с одинаковыми айди не получится.
func (r *Repository) InsertOrder(ctx context.Context, order models.Order) error {

	log := log.Logger.With().Str("func", "Repository.InsertOrder").Caller().Logger()

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
	// т.к. не предполагается конкурентных запросов на чтение --
	// только INSERT'ы. единственная проблема, с которой можем столкнуться,
	// "lost update". Read Uncommited предотвращает lost update.
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, insertOrdersAndDeliverySQL, insertOrderAndDeliveryArgs...); err != nil {
		err = r.handleError(err, log)
		return err
	}

	if _, err := tx.Exec(ctx, insertPaymentSQL, insertPaymentArgs...); err != nil {
		err = r.handleError(err, log)
		return err
	}
	if _, err := tx.Exec(ctx, insertItemQuery, insertItemArgs...); err != nil {
		err = r.handleError(err, log)
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func (*Repository) handleError(err error, log zerolog.Logger) error {
	if pgerr, ok := err.(*pgconn.PgError); ok {
		if pgerrcode.IsIntegrityConstraintViolation(pgerr.Code) {
			log.Debug().Msg("unique_violation")
			return service.ErrDuplicatesNotAllowed
		}
		log.Error().Str("table", pgerr.TableName).Str("column", pgerr.ColumnName).Str("detail", pgerr.Detail).Str("hint", pgerr.Hint).Send()
	} else {
		log.Error().Msg("unknown error")
	}
	return err
}

func (*Repository) paymentArgs(payment models.Payment) []interface{} {
	args := []interface{}{
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDT.Time,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	}

	return args
}

func (*Repository) orderAndDeliveryArgs(order models.Order) []interface{} {
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
