package account

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/external/fxratemanager"
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
	destType := constant.TRANSFER_DESTINATION_UNKNOWN
	if request.Destination.WalletID != "" {
		destType = constant.TRANSFER_DESTINATION_WALLET
	} else if request.Destination.Upi != "" {
		destType = constant.TRANSFER_DESTINATION_UPI
	} else if request.Destination.Account != "" {
		destType = constant.TRANSFER_DESTINATION_BANK
	} else if request.Destination.RecipientDetail != nil {
		// if recipient detail contains cash pickup info
		if _, ok := request.Destination.RecipientDetail["pickup_location"]; ok {
			destType = constant.TRANSFER_DESTINATION_CASH_PICKUP
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
}
