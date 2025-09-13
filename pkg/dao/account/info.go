package account

import (
	"context"
	"database/sql"
	"strings"

	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"gorm.io/gorm"
)

func (obj *accountSt) GetUserAccountsByUserId(ctx context.Context, userId string) ([]dto.DBUserAccount, error) {
	query := `
		select 
			serial_id, user_id, type, bank_name, account_number, ifsc, linked_via, wallet_name, wallet_id, upi_id, created_at, updated_at
		from user_accounts
			where user_id = ?
			and deleted_at is null;
	`

	var (
		response []dto.DBUserAccount = make([]dto.DBUserAccount, 0)
	)

	rows, err := obj.psql.WithContext(ctx).Raw(query, userId).Rows()
	if err != nil {
		if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
			// no user account found
			return nil, nil
		}
		return nil, err
	}

	// scan the data
	for rows.Next() {
		var account dto.DBUserAccount
		if err := rows.Scan(&account.Sid, &account.UserId, &account.Type, &account.BankName, &account.AccountNumber, &account.Ifsc, &account.LinkedVia, &account.WalletName, &account.WalletID, &account.UpiID, &account.CreatedAt, &account.UpdatedAt); err != nil {
			if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
				// no user account found
				return nil, nil
			}
			return nil, err
		}
		response = append(response, account)
	}

	return response, nil
}

func (obj *accountSt) SaveUserAccounts(ctx context.Context, userId string, accounts []dto.UserAccount) error {
	query := `
		INSERT INTO user_accounts (
			user_id, type, bank_name, account_number, 
			ifsc, linked_via, wallet_name, wallet_id, upi_id 
		) VALUES 
	`

	valueQuery := make([]string, 0)
	values := make([]interface{}, 0)

	for _, account := range accounts {
		valueQuery = append(valueQuery, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		values = append(values, userId, account.Type, account.BankName, account.AccountNumber, account.Ifsc, account.LinkedVia, account.WalletName, account.WalletID, account.UpiID)
	}
	query += strings.Join(valueQuery, ", ")
	return obj.psql.WithContext(ctx).Exec(query, values...).Error
}
