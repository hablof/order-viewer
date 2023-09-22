package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/hablof/order-viewer/internal/models"
)

const (
	orderColumns    string = "order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard"
	deliveryColumns string = "name, phone, zip, city, address, region, email"
	paymentColumns  string = "transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee"
)

// Иметь prepared stmt для ReadAll не имеет смысла -- выполняется один раз на старте приложения

func (r *Repository) ReadAll(ctx context.Context) (map[string]models.Order, error) {
	selectOrderDeliveryPaymentQuery, _, err := r.initQuery.
		Select(orderColumns, deliveryColumns, paymentColumns).
		From("orders").
		InnerJoin("delivery USING(delivery_id)").
		InnerJoin("payment ON orders.order_uid = payment.transaction").
		ToSql()

	if err != nil {
		return nil, err
	}

	selectItemQuery, _, err := r.initQuery.Select(itemColumns).From("item").ToSql()
	if err != nil {
		return nil, err
	}

	ctx, cf := context.WithTimeout(ctx, queryTimeout)
	defer cf()

	orderRows, err := r.conn.Query(ctx, selectOrderDeliveryPaymentQuery)
	if err != nil {
		return nil, err
	}
	defer orderRows.Close()

	orders := make(map[string]models.Order) // order_uid -> order

	for orderRows.Next() {
		newOrder, err := scanOrder(orderRows)
		if err != nil {
			return nil, err
		}

		orders[newOrder.OrderUID] = newOrder
	}
	orderRows.Close()

	itemRows, err := r.conn.Query(ctx, selectItemQuery)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		newItem, err := scanItem(itemRows)
		if err != nil {
			return nil, err
		}

		parentOrder := orders[newItem.OrderUID]
		parentOrder.Items = append(parentOrder.Items, newItem)
		orders[newItem.OrderUID] = parentOrder
	}

	return orders, nil
}

// Хотел использовать github.com/georgysavva/scany/v2,
// но он не умеет находить поля во вложенных структурах...
func scanOrder(orderRows pgx.Rows) (models.Order, error) {

	newOrder := models.Order{}
	err := orderRows.Scan(
		&newOrder.OrderUID,
		&newOrder.TrackNumber,
		&newOrder.Entry,
		&newOrder.Locale,
		&newOrder.InternalSignature,
		&newOrder.CustomerID,
		&newOrder.DeliveryService,
		&newOrder.ShardKey,
		&newOrder.SMID,
		&newOrder.DateCreated,
		&newOrder.OofShard,

		&newOrder.Delivery.Name,
		&newOrder.Delivery.Phone,
		&newOrder.Delivery.Zip,
		&newOrder.Delivery.City,
		&newOrder.Delivery.Address,
		&newOrder.Delivery.Region,
		&newOrder.Delivery.Email,

		&newOrder.Payment.Transaction,
		&newOrder.Payment.RequestID,
		&newOrder.Payment.Currency,
		&newOrder.Payment.Provider,
		&newOrder.Payment.Amount,
		&newOrder.Payment.PaymentDT.Time,
		&newOrder.Payment.Bank,
		&newOrder.Payment.DeliveryCost,
		&newOrder.Payment.GoodsTotal,
		&newOrder.Payment.CustomFee,
	)

	return newOrder, err
}

func scanItem(itemRows pgx.Rows) (models.Item, error) {

	newItem := models.Item{}
	err := itemRows.Scan(
		&newItem.OrderUID,
		&newItem.ChrtID,
		&newItem.TrackNumber,
		&newItem.Price,
		&newItem.RID,
		&newItem.Name,
		&newItem.Sale,
		&newItem.Size,
		&newItem.TotalPrice,
		&newItem.NMID,
		&newItem.Brand,
		&newItem.Status,
	)

	return newItem, err
}
