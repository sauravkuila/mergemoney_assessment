package fxratemanager

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"github.com/sauravkuila/mergemoney_assessment/pkg/utils"
	"go.uber.org/zap"
)

// returns the fx rate from sourceCurrency to destinationCurrency along with the date of the rate
func GetFxRate(ctx context.Context, sourceCurrency, destinationCurrency string, utilObj utils.UtilsItf) (float64, string, error) {
	var dataObj RateResponse

	// prepare query params
	date := time.Now().Format("2006-01-02")
	queryParams := map[string]string{
		"base":       sourceCurrency,
		"quote":      destinationCurrency,
		"data_type":  "general_currency_pair",
		"start_date": date,
		"end_date":   date,
	}
	// Make the GET request to fetch the FX rate
	data, code, err := utilObj.InvokeRequest(ctx, http.MethodGet, FX_RATE_VENDOR_BASE_URL+"/cc-api/currencies", nil, nil, nil, queryParams, nil, nil)
	if err != nil || code != http.StatusOK {
		logger.Log(ctx).Error("error in fetching fx rate from vendor", zap.Error(err), zap.Int("status_code", code))
		return 0, date, err
	}

	err = json.Unmarshal(data, &dataObj)
	if err != nil {
		logger.Log(ctx).Error("error in unmarshalling fx rate response", zap.Error(err))
		return 0, date, err
	}

	if len(dataObj.Response) == 0 {
		logger.Log(ctx).Error("no fx rate data found", zap.String("base_currency", sourceCurrency), zap.String("quote_currency", destinationCurrency))
		return 0, date, nil
	}

	// parse the rate, ignore the error
	rate, _ := strconv.ParseFloat(dataObj.Response[0].AverageBid, 64)
	return rate, date, nil
}
