package login

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"github.com/sauravkuila/mergemoney_assessment/pkg/middleware"
	"go.uber.org/zap"
)

func (obj *loginSt) GenerateOTP(c *gin.Context) {
	var (
		request  dto.OTPRequest
		response dto.OTPResponse
	)
	if err := c.BindQuery(&request); err != nil {
		logger.Log(c).Error("failed to bind request", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Generate OTP logic here. call otp service or generate here
	// Save OTP and request ID in database or cache for later verification

	// For demonstration, let's assume we generate a dummy OTP and request ID
	response.Data = &dto.OTPData{
		OTP:          "123456",        // This should be generated dynamically
		OTPRequestId: "req_123456789", // This should be generated dynamically
	}
	response.Status = true
	response.Description = "OTP generated successfully"

	c.JSON(http.StatusOK, response)
}

func (obj *loginSt) VerifyOTP(c *gin.Context) {
	var (
		request  dto.VerifyOTPRequest
		response dto.VerifyOTPResponse
	)
	if err := c.BindJSON(&request); err != nil {
		logger.Log(c).Error("failed to bind request", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// check if OTPRequestId is in our database. rate limit this endpoint based on request id and ip

	// For demonstration, let's assume we verify the OTP request id as req_123456789
	if request.OTPRequestId != "req_123456789" {
		response.Error = append(response.Error, "Invalid OTP request ID")
		response.Description = "Invalid OTP request ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Verify OTP logic here
	// For demonstration, let's assume the OTP is always and is 123456
	if request.OTP != "123456" {
		response.Error = append(response.Error, "Invalid OTP")
		response.Description = "Invalid OTP"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	logger.Log(c).Debug("OTP verified successfully", zap.String("mobile", request.Mobile), zap.String("countryCode", request.CountryCode))

	//fetch user info from database. change this to the OTP req id tracking redis or db
	userRef, err := obj.DB.GetUserFromMobile(c, request.Mobile, request.CountryCode)
	if err != nil {
		logger.Log(c).Error("failed to fetch user data", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Failed to fetch user data"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	if userRef == nil {
		response.Error = append(response.Error, "User not found")
		response.Description = "User not found"
		c.JSON(http.StatusNotFound, response)
		return
	}

	//generate an access and refresh token
	jwtSecret := config.GetConfig().GetString("auth.key")
	claims := map[string]interface{}{
		"user_id": userRef.UserId.String,
		"mobile":  request.Mobile,
		"country": request.CountryCode,
		"type":    "1fa",
		"iat":     time.Now().Unix(),
		"scope":   "non-transactional",
		"token":   "access",
	}
	oneFATokenTimeout := time.Duration(config.GetConfig().GetInt64("auth.key1fa_token_timeout") * int64(time.Second))
	token, err := middleware.GenerateJWT(jwtSecret, claims, time.Now().Add(oneFATokenTimeout))
	if err != nil {
		logger.Log(c).Error("failed to generate JWT token", zap.Error(err))
		response.Error = append(response.Error, "Failed to generate JWT token")
		response.Description = "Internal server error"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	//generate a refresh token
	refreshClaims := map[string]interface{}{
		"user_id": userRef.UserId.String,
		"mobile":  request.Mobile,
		"country": request.CountryCode,
		"type":    "1fa",
		"iat":     time.Now().Unix(),
		"scope":   "non-transactional",
		"token":   "refresh",
	}
	oneFARefreshTimeout := time.Duration(config.GetConfig().GetInt64("auth.key1fa_refresh_timeout") * int64(time.Second))
	refreshToken, err := middleware.GenerateJWT(jwtSecret, refreshClaims, time.Now().Add(oneFARefreshTimeout))
	if err != nil {
		logger.Log(c).Error("failed to generate JWT token", zap.Error(err))
		response.Error = append(response.Error, "Failed to generate JWT token")
		response.Description = "Internal server error"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	response.Description = "OTP verified successfully"
	response.Data = &dto.VerifyOTPData{
		OneFAAccessToken:  token,
		OneFARefreshToken: refreshToken,
		IsMpinPresent:     userRef.Mpin.Valid,
	}

	c.JSON(http.StatusOK, response)
}
