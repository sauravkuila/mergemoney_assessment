package dto

type OTPRequest struct {
	Mobile      string `form:"mobile" binding:"required"`
	CountryCode string `form:"country_code" binding:"required"`
}

type OTPResponse struct {
	Data *OTPData `json:"data,omitempty"`
	CommonResponse
}

type OTPData struct {
	OTP          string `json:"otp"`
	OTPRequestId string `json:"otp_request_id"`
}

type VerifyOTPRequest struct {
	Mobile       string `json:"mobile" binding:"required"`
	CountryCode  string `json:"country_code" binding:"required"`
	OTP          string `json:"otp" binding:"required"`
	OTPRequestId string `json:"otp_request_id" binding:"required"`
}

type VerifyOTPResponse struct {
	Data *VerifyOTPData `json:"data,omitempty"`
	CommonResponse
}

type VerifyOTPData struct {
	OneFAAccessToken  string `json:"one_fa_access_token"`
	OneFARefreshToken string `json:"one_fa_refresh_token"`
	IsMpinPresent     bool   `json:"is_mpin_present"`
}
