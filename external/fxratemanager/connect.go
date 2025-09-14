package fxratemanager

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"github.com/sauravkuila/mergemoney_assessment/pkg/utils"
	"go.uber.org/zap"
)

// returns the fx rate from sourceCurrency to destinationCurrency along with the date of the rate
func GetFxRate(ctx context.Context, sourceCurrency, destinationCurrency string, utilObj utils.UtilsItf) (float64, time.Time, error) {
	var dataObj RateResponse

	// prepare query params
	startDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02") // yesterday
	endDate := time.Now().Format("2006-01-02")
	rateDate := time.Now()
	queryParams := map[string]string{
		"base":       sourceCurrency,
		"quote":      destinationCurrency,
		"data_type":  "general_currency_pair",
		"start_date": startDate,
		"end_date":   endDate,
	}
	u, err := url.Parse(FX_RATE_VENDOR_BASE_URL)
	host := "api.currencycloud.com"
	if err == nil && u.Host != "" {
		host = u.Host
	}

	headers := map[string]string{
		"Host":               host,
		"sec-ch-ua-platform": `"macOS"`,
		"Referer":            "https://www.oanda.com/",
		"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
		"Accept":             "application/json, text/plain, */*",
		"sec-ch-ua":          `"Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"`,
		"sec-ch-ua-mobile":   "?0",
	}
	// Make the GET request to fetch the FX rate
	data, code, err := utilObj.InvokeRequest(ctx, http.MethodGet, FX_RATE_VENDOR_BASE_URL+"/cc-api/currencies", nil, headers, nil, queryParams, nil, nil)
	if err != nil || code != http.StatusOK {
		logger.Log(ctx).Error("error in fetching fx rate from vendor", zap.Error(err), zap.Int("status_code", code), zap.String("data response", string(data)))
		return 0, rateDate, err
	}

	err = json.Unmarshal(data, &dataObj)
	if err != nil {
		logger.Log(ctx).Error("error in unmarshalling fx rate response", zap.Error(err))
		return 0, rateDate, err
	}

	if len(dataObj.Response) == 0 {
		logger.Log(ctx).Error("no fx rate data found", zap.String("base_currency", sourceCurrency), zap.String("quote_currency", destinationCurrency))
		return 0, rateDate, nil
	}

	// parse the rate, ignore the error
	rate, _ := strconv.ParseFloat(dataObj.Response[0].AverageBid, 64)
	rateDate, err = time.Parse(time.RFC3339, dataObj.Response[0].CloseTime)
	if err != nil {
		logger.Log(ctx).Error("error parsing close_time", zap.Error(err), zap.String("close_time", dataObj.Response[0].CloseTime))
	}
	return rate, rateDate, nil
}
