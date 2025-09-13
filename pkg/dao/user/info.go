package user

import (
	"context"

	"gorm.io/gorm"
)

func (obj *userSt) SetUserMPIN(ctx context.Context, userId string, mpin string) error {
	query := `
		UPDATE user_ref
		SET user_mpin = ?
		WHERE user_id = ? AND deleted_at IS NULL;
	`

	result := obj.psql.WithContext(ctx).Exec(query, mpin, userId)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (obj *userSt) ResetUserMPIN(c context.Context, mobile string, countryCode string) error {
	query := `
		UPDATE user_ref
		SET user_mpin = NULL
		WHERE mobile = ? AND country_code = ? AND deleted_at IS NULL;
	`

	result := obj.psql.WithContext(c).Exec(query, mobile, countryCode)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
