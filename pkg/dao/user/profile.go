package user

// func (obj *userSt) UpdateUsersData(ctx context.Context, data dto.DbBrokerUserData) error {
// 	query := `
// 		UPDATE user_broker_info
// 	`

// 	var (
// 		setQueries []string
// 		whereQuery string
// 		args       []interface{}
// 	)
// 	// add set queries
// 	if len(data.Metadata) != 0 {
// 		setQueries = append(setQueries, "SET metadata = ?")
// 	}
// 	if data.BrokerUserIdentifier.Valid {
// 		setQueries = append(setQueries, "SET broker_user_identifier = ?")
// 	}
// 	args = append(args, data.Metadata)
// 	if data.BrokerUserIdentifier.Valid {
// 		whereQuery = "WHERE broker_user_identifier = ?"
// 		args = append(args, data.BrokerUserIdentifier)
// 	} else if data.UserId.Valid && data.BrokerId.Valid {
// 		whereQuery = "WHERE (user_id = ? AND broker_id = ?)"
// 		args = append(args, data.UserId, data.BrokerId)
// 	} else {
// 		return fmt.Errorf("invalid arguments. require userId+brokerId OR brokerUserIdentifier")
// 	}

// 	// prepare query
// 	query += strings.Join(setQueries, ",") + " " + whereQuery

// 	result := obj.psql.WithContext(ctx).Exec(query, args...)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	if result.RowsAffected == 0 {
// 		return fmt.Errorf("no rows affected")
// 	}
// 	return nil
// }

// func (obj *userSt) GetUserData(ctx context.Context, data dto.DbBrokerUserData) (*dto.DbBrokerUserData, error) {
// 	// query := `
// 	// 	SELECT metadata, broker_user_identifier
// 	// 	FROM user_broker_info
// 	// `
// 	query := `
// 		SELECT ubi.metadata, ubi.broker_user_identifier, uwt.token
// 		FROM user_broker_info ubi
// 		LEFT OUTER JOIN user_wss_token uwt ON ubi.user_id = uwt.user_id
// 	`

// 	var (
// 		queries []string
// 		args    []interface{}
// 	)
// 	queries = append(queries, query)
// 	// args = append(args, data.Metadata)
// 	if data.BrokerUserIdentifier.Valid {
// 		queries = append(queries, "WHERE ubi.broker_user_identifier = ?")
// 		args = append(args, data.BrokerUserIdentifier)
// 	} else if data.UserId.Valid && data.BrokerId.Valid {
// 		queries = append(queries, "WHERE (ubi.user_id = ? AND ubi.broker_id = ?)")
// 		args = append(args, data.UserId, data.BrokerId)
// 	} else {
// 		return nil, fmt.Errorf("invalid arguments. require userId+brokerId OR brokerUserIdentifier")
// 	}

// 	queries = append(queries, "AND ubi.deleted_at IS NULL ORDER BY uwt.created_at DESC LIMIT 1;")
// 	// prepare the executable query
// 	query = strings.Join(queries, " ")

// 	var (
// 		metadataBytes sql.NullString
// 		response      dto.DbBrokerUserData
// 	)

// 	row := obj.psql.WithContext(ctx).Raw(query, args...).Row()
// 	if row.Err() != nil {
// 		return nil, row.Err()
// 	}

// 	// scan the data
// 	if err := row.Scan(&metadataBytes, &response.BrokerUserIdentifier, &response.WssToken); err != nil {
// 		return nil, err
// 	}

// 	if !metadataBytes.Valid {
// 		return nil, sql.ErrNoRows
// 	}
// 	response.Metadata = []byte(metadataBytes.String)

// 	return &response, nil
// }
