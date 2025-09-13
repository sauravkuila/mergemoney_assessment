package user

import (
	"context"
	"database/sql"
	"log"

	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"gorm.io/gorm"
)

func (obj *userSt) GetUserFromMobile(ctx context.Context, mobile string, countryCode string) (*dto.DBUserRef, error) {
	query := `
		select 
			user_id, user_name, mobile, country_code, user_role, user_mpin 
		from user_ref
			where mobile = ? and country_code like ?
			and deleted_at is null;
	`

	var (
		response dto.DBUserRef
	)

	row := obj.psql.WithContext(ctx).Raw(query, mobile, countryCode).Row()
	if row.Err() != nil {
		log.Println(" row error", row.Err())
		if row.Err() == sql.ErrNoRows || row.Err() == gorm.ErrRecordNotFound {
			// no user found
			return nil, nil
		}
		return nil, row.Err()
	}

	// scan the data
	if err := row.Scan(&response.UserId, &response.UserName, &response.Mobile, &response.CountryCode, &response.UserRole, &response.Mpin); err != nil {
		if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
			// no user found
			return nil, nil
		}
		return nil, err
	}

	return &response, nil
}
