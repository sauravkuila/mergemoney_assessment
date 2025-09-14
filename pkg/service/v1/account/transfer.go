package account

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/external/fxratemanager"
	"github.com/sauravkuila/mergemoney_assessment/external/paymentprovider"
	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/constant"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"go.uber.org/zap"
)

// Transfer handler
func (obj *accountSt) Transfer(c *gin.Context) {
	var (
		request  dto.TransferRequest
		response dto.TransferResponse
	)
	if err := c.BindJSON(&request); err != nil {
		logger.Log(c).Error("error in binding transfer request", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	logger.Log(c).Info("transfer request", zap.Any("request", request))

	err := obj.validateTransferRequest(c, request)
	if err != nil {
		logger.Log(c).Error("error in validating transfer request", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// determine destination type
	destType := constant.TRANSFER_TYPE_UNKNOWN
	if request.Destination.WalletID != "" {
		destType = constant.TRANSFER_TYPE_WALLET
	} else if request.Destination.Upi != "" {
		destType = constant.TRANSFER_TYPE_UPI
	} else if request.Destination.Account != "" {
		destType = constant.TRANSFER_TYPE_BANK
	} else if request.Destination.RecipientDetail != nil {
		// if recipient detail contains cash pickup info
		if _, ok := request.Destination.RecipientDetail["pickup_location"]; ok {
			destType = constant.TRANSFER_TYPE_CASH_PICKUP
		}
	} else {
		logger.Log(c).Error("invalid destination details", zap.Any("destination", request.Destination))
		response.Error = append(response.Error, "invalid destination details")
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// call payment provider api to initiate transfer - simulate success response for now
	rate, rateDate, err := fxratemanager.GetFxRate(c, request.Source.Currency, request.Destination.Currency, obj.utilsObj)
	if err != nil {
		logger.Log(c).Error("error in fetching fx rate", zap.Error(err))
		response.Error = append(response.Error, "error in fetching fx rate")
		response.Description = "Unable to initiate transfer"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	logger.Log(c).Info("fx rate fetched", zap.Float64("rate", rate), zap.Time("rate_date", rateDate))
	if rate == 0 {
		rate = 1.0 // assume same currency if rate not found
	}

	// // generate an idempotent transfer id using utils Identifier
	id := obj.utilsObj.GetUniqueId("ordr")

	// create an order in orders table
	order := dto.DBOrder{
		OrderID:             sql.NullString{String: id, Valid: true},
		UserID:              sql.NullString{String: c.GetString(config.USERID), Valid: true},
		SourceSID:           sql.NullInt64{Int64: request.Source.Sid, Valid: true},
		SourceCurrency:      sql.NullString{String: request.Source.Currency, Valid: true},
		DestinationCurrency: sql.NullString{String: request.Destination.Currency, Valid: true},
		SourceAmount:        sql.NullFloat64{Float64: request.Source.Amount, Valid: true},
		DestinationAmount:   sql.NullFloat64{Float64: request.Source.Amount * rate, Valid: true},
		ConversionRate:      sql.NullFloat64{Float64: rate, Valid: true},
		ConversionRateDate:  sql.NullTime{Time: rateDate, Valid: true},
		OrderStatus:         sql.NullString{String: constant.ORDER_CREATED, Valid: true},
		Remark:              sql.NullString{String: "", Valid: false},
	}
	orderDestination := dto.DBOrderDestination{
		OrderID:           sql.NullString{String: id, Valid: true},
		DestinationType:   sql.NullString{String: destType, Valid: true},
		WalletID:          sql.NullString{String: request.Destination.WalletID, Valid: request.Destination.WalletID != ""},
		UPIID:             sql.NullString{String: request.Destination.Upi, Valid: request.Destination.Upi != ""},
		BankAccountNumber: sql.NullString{String: request.Destination.Account, Valid: request.Destination.Account != ""},
		IFSCCode:          sql.NullString{String: request.Destination.SwiftCode, Valid: request.Destination.SwiftCode != ""},
	}

	// save the order in DB
	if err := obj.DB.SaveOrder(c, order, orderDestination); err != nil {
		logger.Log(c).Error("error in saving transfer to DB", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Unable to save transfer"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	response.Description = "Transfer created"
	response.Data = &dto.TransferData{
		TransferID:          id,
		SourceCurrency:      request.Source.Currency,
		DestinationCurrency: request.Destination.Currency,
		SourceAmount:        request.Source.Amount,
		DestinationAmount:   request.Source.Amount * rate,
		ConversionRate:      rate,
		ConversionRateDate:  rateDate.Format(time.DateTime),
		DestinationType:     destType,
	}

	c.JSON(http.StatusOK, response)
}

func (obj *accountSt) validateTransferRequest(c *gin.Context, transfer dto.TransferRequest) error {
	// validate source fields
	if transfer.Source.Sid <= 0 {
		logger.Log(c).Error("invalid source sid", zap.Int64("sid", transfer.Source.Sid))
		return errors.New("invalid source sid")
	} else {
		// check if source sid belongs to user
		userId := c.GetString(config.USERID)
		accs, err := obj.DB.GetUserAccountsByUserId(c, userId)
		if err != nil || len(accs) == 0 {
			logger.Log(c).Error("error in fetching user accounts from DB", zap.Error(err), zap.String("userId", userId))
			return errors.New("invalid source sid")
		}
		found := false
		for _, a := range accs {
			if a.Sid.Int64 == transfer.Source.Sid {
				found = true
				break
			}
		}
		if !found {
			return errors.New("invalid source sid")
		}
	}

	if transfer.Source.Amount <= 0 {
		logger.Log(c).Error("invalid source amount", zap.Float64("amount", transfer.Source.Amount))
		return errors.New("invalid source amount")
	}

	cCode := transfer.Source.Currency
	if len(cCode) != 3 {
		logger.Log(c).Error("invalid currency code", zap.String("currency", cCode))
		return errors.New("invalid currency code")
	}
	for _, ch := range cCode {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')) {
			logger.Log(c).Error("invalid currency code", zap.String("currency", cCode))
			return errors.New("invalid currency code")
		}
	}

	return nil
}

func (obj *accountSt) TransferConfirm(c *gin.Context) {
	var (
		request  dto.TransferConfirmRequest
		response dto.TransferConfirmResponse
	)

	if err := c.BindJSON(&request); err != nil {
		logger.Log(c).Error("error in binding transfer confirm request", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	logger.Log(c).Info("transfer confirm request", zap.Any("request", request))

	// fetch the order from DB and check if it belongs to user and is in pending state
	order, orderDest, err := obj.DB.GetOrderById(c, request.TransferID, c.GetString(config.USERID))
	if err != nil {
		logger.Log(c).Error("error in fetching order from DB", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid transfer ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// fetch source detail from DB using sid
	userAccount, err := obj.DB.GetUserAccountsBySid(c, c.GetString(config.USERID), order.SourceSID.Int64)
	if err != nil {
		logger.Log(c).Error("error in fetching user accounts from DB", zap.Error(err), zap.String("userId", c.GetString(config.USERID)))
		response.Error = append(response.Error, "error in fetching source account")
		response.Description = "Unable to process transfer"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	//check status of order
	if order.OrderStatus.String != constant.ORDER_CREATED {
		logger.Log(c).Error("order not in created state", zap.String("order_id", order.OrderID.String), zap.String("status", order.OrderStatus.String))
		response.Error = append(response.Error, "order not in created state")
		response.Description = "Invalid transfer state"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// execution reaches here is order is valid and in created state.
	// order can either be cancelled or send to provider for processing based on action
	switch request.Action {
	case "cancel":
		// update order status to cancelled
		err := obj.DB.UpdateOrderStatus(c, order.OrderID.String, constant.ORDER_FAILED, "cancelled by user")
		if err != nil {
			logger.Log(c).Error("error in updating order status", zap.Error(err))
			response.Error = append(response.Error, err.Error())
			response.Description = "Failed to cancel transfer"
			c.JSON(http.StatusInternalServerError, response)
			return
		}
	case "confirm":
		// call payment provider api to process the transfer
		logger.Log(c).Info("calling payment provider api to process transfer", zap.String("order_id", order.OrderID.String), zap.Any("destination", orderDest))

		payDetail := paymentprovider.PaymentDetails{}

		// fill source details
		switch userAccount.Type.String {
		case constant.TRANSFER_TYPE_WALLET:
			payDetail.SourceDetail.Type = paymentprovider.PayTypeWallet
			// fill source wallet id
			payDetail.SourceDetail.WalletID = userAccount.WalletID.String
		case constant.TRANSFER_TYPE_UPI:
			payDetail.SourceDetail.Type = paymentprovider.PayTypeUPI
			// fill upi id
			payDetail.SourceDetail.UPIID = userAccount.UpiID.String
		case constant.TRANSFER_TYPE_BANK:
			payDetail.SourceDetail.Type = paymentprovider.PayTypeNetBanking
			// fill account and swift code
			payDetail.SourceDetail.AccountNumber = userAccount.AccountNumber.String
			payDetail.SourceDetail.SwiftCode = userAccount.Ifsc.String
		default:
			payDetail.SourceDetail.Type = paymentprovider.PayTypeCash
		}

		// fill destination details
		switch orderDest.DestinationType.String {
		case constant.TRANSFER_TYPE_WALLET:
			payDetail.DestinationDetail.Type = paymentprovider.PayTypeWallet
			// fill source wallet id
			payDetail.DestinationDetail.WalletID = orderDest.WalletID.String
		case constant.TRANSFER_TYPE_UPI:
			payDetail.DestinationDetail.Type = paymentprovider.PayTypeUPI
			// fill upi id
			payDetail.DestinationDetail.UPIID = orderDest.UPIID.String
		case constant.TRANSFER_TYPE_BANK:
			payDetail.DestinationDetail.Type = paymentprovider.PayTypeNetBanking
			// fill account and swift code
			payDetail.DestinationDetail.AccountNumber = orderDest.BankAccountNumber.String
			payDetail.DestinationDetail.SwiftCode = orderDest.IFSCCode.String
		default:
			payDetail.DestinationDetail.Type = paymentprovider.PayTypeCash
		}

		transferId, provider, providerRequest, providerResponse, err := paymentprovider.InitiateTransfer(c, order.DestinationAmount.Float64, order.DestinationCurrency.String, payDetail, obj.utilsObj)
		if err != nil {
			logger.Log(c).Error("error initiating transfer", zap.Error(err))
			response.Error = append(response.Error, err.Error())
			response.Description = "Failed to initiate transfer"
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// update db with provider transfer id and update order status to inprogress
		err = obj.DB.UpdateOrderStatus(c, order.OrderID.String, constant.ORDER_INPROGRESS, "transfer initiated with provider: "+transferId)
		if err != nil {
			logger.Log(c).Error("error in updating order status", zap.Error(err))
			response.Error = append(response.Error, err.Error())
			response.Description = "Failed to process transfer"
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		logger.Log(c).Info("transfer initiated successfully", zap.String("order_id", order.OrderID.String), zap.String("provider_transfer_id", transferId), zap.Any("provider", provider))

		// update the transaction table with provider details
		txn := dto.DBTransaction{
			TransactionID:    sql.NullString{String: obj.utilsObj.GetUniqueId("txn"), Valid: true},
			OrderID:          sql.NullString{String: request.TransferID, Valid: true},
			Provider:         sql.NullString{String: provider, Valid: true},
			ProviderID:       sql.NullString{String: transferId, Valid: true},
			ProviderRequest:  providerRequest,
			ProviderResponse: providerResponse,
			Status:           sql.NullString{String: constant.TXN_INITIATED, Valid: true},
			ErrorMessage:     sql.NullString{String: "", Valid: false},
			RetryCount:       sql.NullInt64{Int64: 0, Valid: true},
			LastRetryAt:      sql.NullTime{Time: time.Time{}, Valid: true},
		}
		if err := obj.DB.SaveTransaction(c, txn); err != nil {
			logger.Log(c).Error("error in saving transaction to DB", zap.Error(err))
			response.Error = append(response.Error, err.Error())
			response.Description = "Failed to process transfer"
			c.JSON(http.StatusInternalServerError, response)
			return
		}

	default:
		logger.Log(c).Error("invalid action", zap.String("action", request.Action))
		response.Error = append(response.Error, "invalid action")
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = true
	if request.Action == "confirm" {
		response.Description = "Transfer confirmed"
	} else {
		response.Description = "Transfer cancelled"
	}
	response.Data = &dto.TransferConfirm{
		TransferID: request.TransferID,
		Status:     response.Description,
	}

	c.JSON(http.StatusOK, response)
}

func (obj *accountSt) TransferStatus(c *gin.Context) {
	var (
		request  dto.TransferStatusRequest
		response dto.TransferStatusResponse
	)

	if err := c.BindUri(&request); err != nil {
		logger.Log(c).Error("error in binding transfer status request", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	logger.Log(c).Info("transfer status request", zap.Any("request", request))

	// fetch the order from DB and check if it belongs to user
	order, _, err := obj.DB.GetOrderById(c, request.TransferID, c.GetString(config.USERID))
	if err != nil {
		logger.Log(c).Error("error in fetching order from DB", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid transfer ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.Status = true
	response.Description = "Transfer status fetched"
	response.Data = &dto.TransferStatus{
		TransferID: order.OrderID.String,
		Status:     order.OrderStatus.String,
	}

	c.JSON(http.StatusOK, response)
}
