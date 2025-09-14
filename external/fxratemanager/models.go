package fxratemanager

type RateResponse struct {
	Response []struct {
		BaseCurrency  string `json:"base_currency"`
		QuoteCurrency string `json:"quote_currency"`
		CloseTime     string `json:"close_time"`
		AverageBid    string `json:"average_bid"`
		AverageAsk    string `json:"average_ask"`
		HighBid       string `json:"high_bid"`
		HighAsk       string `json:"high_ask"`
		LowBid        string `json:"low_bid"`
		LowAsk        string `json:"low_ask"`
	} `json:"response"`
}
