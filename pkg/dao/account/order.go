package account

import (
	"context"

	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
)

func (obj *accountSt) SaveOrder(ctx context.Context, order dto.DBOrder, orderDestination dto.DBOrderDestination) error {
	query := `
		INSERT INTO orders (
			order_id, user_id, source_sid, source_currency, source_amount,
			destination_currency, destination_amount, conversion_rate, conversion_rate_date, order_status
		) VALUES (
		 	?,?,?,?,?,
			?,?,?,?,?
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

	orderTx := tx.WithContext(ctx).Exec(query, order.OrderID, order.UserID, order.SourceSID, order.SourceCurrency, order.SourceAmount, order.DestinationCurrency, order.DestinationAmount, order.ConversionRate, order.ConversionRateDate, order.OrderStatus)
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
