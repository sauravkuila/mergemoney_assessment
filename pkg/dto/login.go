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

type SetMPINRequest struct {
	Mpin string `json:"mpin" binding:"required"`
}

type SetMPINResponse struct {
	Data string `json:"data,omitempty"`
	CommonResponse
}

type ResetMPINRequest struct {
	Mobile      string `json:"mobile" binding:"required"`
	CountryCode string `json:"country_code" binding:"required"`
}

type ResetMPINResponse struct {
	Data bool `json:"data,omitempty"`
	CommonResponse
}

type VerifyMPINRequest struct {
	Mpin string `json:"mpin" binding:"required"`
}

type VerifyMPINResponse struct {
	Data *VerifyMPINData `json:"data,omitempty"`
	CommonResponse
}

type VerifyMPINData struct {
	AccessToken        string `json:"access_token"`
	AccessTokenExpiry  int64  `json:"access_token_expiry"`
	RefreshToken       string `json:"refresh_token"`
	RefreshTokenExpiry int64  `json:"refresh_token_expiry"`
}

type Refresh1FARequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type Refresh1FAAResponse struct {
	Data *Refresh1FA `json:"data,omitempty"`
	CommonResponse
}

type Refresh1FA struct {
	OneFAAccessToken string `json:"one_fa_access_token"`
}

type Refresh2FARequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type Refresh2FAResponse struct {
	Data *Refresh2FA `json:"data,omitempty"`
	CommonResponse
}

type Refresh2FA struct {
	AccessToken string `json:"access_token"`
}
