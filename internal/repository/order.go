package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/underthetreee/L0/internal/model"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (s *PostgresRepository) Store(ctx context.Context, order model.Order) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("transaction begin: %w", err)
	}

	if err := insertOrder(ctx, tx, order); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("transaction rollback: %w", err)
		}
		return fmt.Errorf("transaction process: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit: %w", err)
	}

	return nil
}

func (s *PostgresRepository) GetAll(ctx context.Context) ([]model.Order, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("transaction begin: %w", err)
	}

	orders, err := getOrders(ctx, tx)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, fmt.Errorf("transaction rollback: %w", err)
		}
		return nil, fmt.Errorf("transaction process: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("transaction commit: %w", err)
	}
	return orders, nil

}

func getOrders(ctx context.Context, tx pgx.Tx) ([]model.Order, error) {
	rows, err := tx.Query(ctx, `SELECT
		o.uid, o.track_number, o.entry,
		d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		p.transaction, p.request_id, p.currency, p.provider, p.amount,
		p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
		o.locale, o.internal_signature, o.custom_id, o.delivery_service,
		o.shard_key, o.sm_id, o.date_created, o.oof_shard
	FROM orders o
	JOIN deliveries d ON o.uid=d.order_uid
	JOIN payments p ON o.uid=p.order_uid`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.UID, &order.TrackNumber, &order.Entry,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
			&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
			&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT,
			&order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee,
			&order.Locale, &order.InternalSignature, &order.CustomID, &order.DeliveryService,
			&order.ShardKey, &order.SMID, &order.DateCreated, &order.OofShard)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func insertOrder(ctx context.Context, tx pgx.Tx, order model.Order) error {
	_, err := tx.Exec(ctx, `
        INSERT INTO orders (uid, track_number, entry, locale,
			internal_signature, custom_id, delivery_service, shard_key, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.UID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomID, order.DeliveryService, order.ShardKey, order.SMID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
	    INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
	    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.UID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
	    INSERT INTO payments (order_uid, transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.UID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(ctx, `
	        INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale,
				size, total_price, nm_id, brand, status)
	        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.UID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}
	return nil
}
