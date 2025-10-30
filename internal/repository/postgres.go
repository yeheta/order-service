package repository

import (
	"database/sql"
	"fmt"
	"order-service/internal/domain"

	_ "github.com/lib/pq"
)

type OrderRepository interface {
	Save(order *domain.Order) error
	GetByUID(uid string) (*domain.Order, error)
	GetAll() ([]domain.Order, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(host string, port int, user, password, dbname string) (*PostgresRepository, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Save(order *domain.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Сохранение delivery
	var deliveryID int
	err = tx.QueryRow(`
		INSERT INTO delivery (name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region,
		order.Delivery.Email).Scan(&deliveryID)
	if err != nil {
		return err
	}

	// Сохранение payment
	_, err = tx.Exec(`
		INSERT INTO payment (transaction, request_id, currency, provider, amount, 
		                   payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (transaction) DO NOTHING
	`, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal,
		order.Payment.CustomFee)
	if err != nil {
		return err
	}

	// Сохранение order
	_, err = tx.Exec(`
		INSERT INTO orders (order_uid, track_number, entry, delivery_id, 
		                   payment_transaction, locale, internal_signature, 
		                   customer_id, delivery_service, shardkey, sm_id, 
		                   date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (order_uid) DO NOTHING
	`, order.OrderUID, order.TrackNumber, order.Entry, deliveryID,
		order.Payment.Transaction, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID,
		order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}

	// Сохранение items
	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (chrt_id, track_number, price, rid, name, sale, 
			                   size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT (id) DO NOTHING
		`, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetByUID(uid string) (*domain.Order, error) {
	query := `
		SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, 
		       o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
		       d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		       p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, 
		       p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		JOIN delivery d ON o.delivery_id = d.id
		JOIN payment p ON o.payment_transaction = p.transaction
		WHERE o.order_uid = $1
	`

	var order domain.Order
	var delivery domain.Delivery
	var payment domain.Payment

	err := r.db.QueryRow(query, uid).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
		&delivery.Address, &delivery.Region, &delivery.Email,
		&payment.Transaction, &payment.RequestID, &payment.Currency,
		&payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank,
		&payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}

	items, err := r.getItemsByOrderUID(uid)
	if err != nil {
		return nil, err
	}

	order.Delivery = delivery
	order.Payment = payment
	order.Items = items

	return &order, nil
}

func (r *PostgresRepository) GetAll() ([]domain.Order, error) {
	query := `
		SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, 
		       o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
		       d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		       p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, 
		       p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		JOIN delivery d ON o.delivery_id = d.id
		JOIN payment p ON o.payment_transaction = p.transaction
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var delivery domain.Delivery
		var payment domain.Payment

		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
			&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
			&delivery.Address, &delivery.Region, &delivery.Email,
			&payment.Transaction, &payment.RequestID, &payment.Currency,
			&payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank,
			&payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
		)
		if err != nil {
			return nil, err
		}

		items, err := r.getItemsByOrderUID(order.OrderUID)
		if err != nil {
			return nil, err
		}

		order.Delivery = delivery
		order.Payment = payment
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PostgresRepository) getItemsByOrderUID(orderUID string) ([]domain.Item, error) {
    query := `
        SELECT i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, 
               i.total_price, i.nm_id, i.brand, i.status
        FROM items i
        JOIN orders o ON i.track_number = o.track_number
        WHERE o.order_uid = $1
    `

    rows, err := r.db.Query(query, orderUID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []domain.Item
    for rows.Next() {
        var item domain.Item
        err := rows.Scan(
            &item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
            &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
        )
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }

    return items, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}