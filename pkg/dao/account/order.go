package account

import (
	"context"
	"database/sql"

	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
)

func (obj *accountSt) SaveOrder(ctx context.Context, order dto.DBOrder, orderDestination dto.DBOrderDestination) error {
	query := `
		INSERT INTO orders (
			order_id, user_id, source_sid, source_currency, source_amount,
			destination_currency, destination_amount, conversion_rate, conversion_rate_date, order_status, remarks
		) VALUES (
		 	?,?,?,?,?,
			?,?,?,?,?,?
		);
	`

	// Begin a transaction
	tx := obj.psql.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Query for inserting into order destination
	orderDestinationQuery := `
		INSERT INTO order_destinations (
			order_id, destination_type, wallet_id, upi_id, bank_account_number, ifsc_code
		) VALUES (
		 	?,?,?,?,?,?
		);
	`

	orderTx := tx.WithContext(ctx).Exec(query, order.OrderID, order.UserID, order.SourceSID, order.SourceCurrency, order.SourceAmount, order.DestinationCurrency, order.DestinationAmount, order.ConversionRate, order.ConversionRateDate, order.OrderStatus, order.Remark)
	if orderTx.Error != nil {
		// Rollback the transaction on error
		tx.Rollback()
		return orderTx.Error
	}

	// insert the order destination
	orderDestinationTx := tx.WithContext(ctx).Exec(orderDestinationQuery, order.OrderID, orderDestination.DestinationType, orderDestination.WalletID, orderDestination.UPIID, orderDestination.BankAccountNumber, orderDestination.IFSCCode)
	if orderDestinationTx.Error != nil {
		// Rollback the transaction on error
		tx.Rollback()
		return orderDestinationTx.Error
	}

	return tx.Commit().Error
}

func (obj *accountSt) GetOrderById(ctx context.Context, orderId string, userId string) (*dto.DBOrder, *dto.DBOrderDestination, error) {
	query := `
		SELECT o.order_id, o.user_id, o.source_sid, o.source_currency, o.source_amount,
			o.destination_currency, o.destination_amount, o.conversion_rate, o.conversion_rate_date, o.order_status, o.remarks,
			o.created_at, o.updated_at, od.destination_type, od.wallet_id, od.upi_id, od.bank_account_number, od.ifsc_code
		FROM orders o
		LEFT OUTER JOIN order_destinations od ON o.order_id = od.order_id
		WHERE o.order_id = ? AND o.user_id = ?;
	`
	var (
		order            dto.DBOrder
		orderDestination dto.DBOrderDestination
	)
	rows, err := obj.psql.WithContext(ctx).Raw(query, orderId, userId).Rows()
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&order.OrderID, &order.UserID, &order.SourceSID, &order.SourceCurrency, &order.SourceAmount,
			&order.DestinationCurrency, &order.DestinationAmount, &order.ConversionRate, &order.ConversionRateDate, &order.OrderStatus, &order.Remark, &order.CreatedAt, &order.UpdatedAt,
			&orderDestination.DestinationType, &orderDestination.WalletID, &orderDestination.UPIID, &orderDestination.BankAccountNumber, &orderDestination.IFSCCode,
		); err != nil {
			return nil, nil, err
		}
	}

	if order.OrderID.String == "" {
		// no order found
		return nil, nil, sql.ErrNoRows
	}

	return &order, &orderDestination, nil
}

func (obj *accountSt) UpdateOrderStatus(ctx context.Context, orderId string, status string, remark string) error {
	query := `
		UPDATE orders
		SET order_status = ?, remarks = ?, updated_at = NOW()
		WHERE order_id = ?;
	`
	result := obj.psql.WithContext(ctx).Exec(query, status, remark, orderId)
	return result.Error
}
