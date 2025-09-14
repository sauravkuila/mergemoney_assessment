package account

import (
	"context"
	"database/sql"

	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
)

// Save the transaction
func (obj *accountSt) SaveTransaction(ctx context.Context, transaction dto.DBTransaction) error {
	query := `
		INSERT INTO transactions (
			transaction_id, order_id, provider, provider_id, provider_request, provider_response,
			status, error_message, retry_count, last_retry_at
		) VALUES (
		 	?,?,?,?,?,
			?,?,?,?,?
		);
	`

	row := obj.psql.WithContext(ctx).Raw(query, transaction.TransactionID, transaction.OrderID, transaction.Provider, transaction.ProviderID, transaction.ProviderRequest, transaction.ProviderResponse, transaction.Status, transaction.ErrorMessage, transaction.RetryCount, transaction.LastRetryAt).Row()
	if row.Err() != nil {
		return row.Err()
	}

	if err := row.Scan(&transaction.TransactionID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}

// Update the transaction status
func (obj *accountSt) UpdateTransactionStatus(ctx context.Context, providerId string, transactionId string, status string, remark string) error {
	query := `
		UPDATE transactions
		SET status = ?
		WHERE provider_id = ? AND transaction_id = ?;
	`
	result := obj.psql.WithContext(ctx).Exec(query, status, providerId, transactionId)
	return result.Error
}

// Get the transaction by id
func (obj *accountSt) GetTransactionById(ctx context.Context, transactionId string) (*dto.DBTransaction, error) {
	query := `
		SELECT 
			transaction_id, order_id, provider, provider_id, provider_request, provider_response,
			status, error_message, retry_count, last_retry_at, created_at, updated_at
		FROM transactions
		WHERE transaction_id = ?;
	`
	var (
		transaction dto.DBTransaction
	)
	row := obj.psql.WithContext(ctx).Raw(query, transactionId).Row()
	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := row.Scan(
		&transaction.TransactionID, &transaction.OrderID, &transaction.Provider, &transaction.ProviderID, &transaction.ProviderRequest, &transaction.ProviderResponse,
		&transaction.Status, &transaction.ErrorMessage, &transaction.RetryCount, &transaction.LastRetryAt, &transaction.CreatedAt, &transaction.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// Get the transaction by order id
func (obj *accountSt) GetTransactionByOrderId(ctx context.Context, orderId string) ([]dto.DBTransaction, error) {
	var transactions []dto.DBTransaction
	query := `
		SELECT 
			transaction_id, order_id, provider, provider_id, provider_request, provider_response,
			status, error_message, retry_count, last_retry_at, created_at, updated_at
		FROM transactions
		WHERE order_id = ?;
	`
	rows, err := obj.psql.WithContext(ctx).Raw(query, orderId).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction dto.DBTransaction
		if err := rows.Scan(
			&transaction.TransactionID, &transaction.OrderID, &transaction.Provider, &transaction.ProviderID, &transaction.ProviderRequest, &transaction.ProviderResponse,
			&transaction.Status, &transaction.ErrorMessage, &transaction.RetryCount, &transaction.LastRetryAt, &transaction.CreatedAt, &transaction.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		// no transaction found
		return nil, sql.ErrNoRows
	}

	return transactions, nil
}
