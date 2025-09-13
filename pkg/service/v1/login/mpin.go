package login

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"github.com/sauravkuila/mergemoney_assessment/pkg/middleware"
	"go.uber.org/zap"
)

func (obj *loginSt) SetMPIN(c *gin.Context) {
	var (
		request  dto.SetMPINRequest
		response dto.SetMPINResponse
	)

	// Bind the request
	if err := c.BindJSON(&request); err != nil {
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// update MPIN in the database
	userId := c.GetString(config.USERID)

	if err := obj.DB.SetUserMPIN(c, userId, request.Mpin); err != nil {
		response.Error = append(response.Error, err.Error())
		response.Description = "Failed to set MPIN"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	encodedMPIN := base64.StdEncoding.EncodeToString([]byte(request.Mpin))
	response.Status = true
	response.Description = "MPIN set successfully"
	response.Data = encodedMPIN

	c.JSON(http.StatusOK, response)
}

func (obj *loginSt) ResetMPIN(c *gin.Context) {
	var (
		request  dto.ResetMPINRequest
		response dto.ResetMPINResponse
	)

	// Bind the request
	if err := c.BindJSON(&request); err != nil {
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// update MPIN in the database
	err := obj.DB.ResetUserMPIN(c, request.Mobile, request.CountryCode)
	if err != nil {
		response.Error = append(response.Error, err.Error())
		response.Description = "Failed to reset MPIN"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	response.Description = "MPIN reset successfully"
	response.Data = true

	c.JSON(http.StatusOK, response)
}

func (obj *loginSt) VerifyMPIN(c *gin.Context) {
	var (
		request  dto.VerifyMPINRequest
		response dto.VerifyMPINResponse
	)
	if err := c.BindJSON(&request); err != nil {
		logger.Log(c).Error("failed to bind request", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userId := c.GetString(config.USERID)

	//fetch user info from database. change this to the OTP req id tracking redis or db
	userRef, err := obj.DB.GetUserFromUserId(c, userId)
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

	// match if user pin and sent pin match
	// TODO: use bcrypt for hashing
	if userRef.Mpin.Valid && userRef.Mpin.String != request.Mpin {
		response.Error = append(response.Error, "Invalid MPIN")
		response.Description = "Invalid MPIN"
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	// if MPIN is not set, return error
	if !userRef.Mpin.Valid {
		response.Error = append(response.Error, "MPIN not set")
		response.Description = "MPIN not set"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// if MPIN is set, generate access token and refresh token

	//generate an access and refresh token
	jwtSecret := config.GetConfig().GetString("auth.key2fa")
	claims := map[string]interface{}{
		"user_id": userRef.UserId.String,
		"mobile":  userRef.Mobile.String,
		"country": userRef.CountryCode.String,
		"type":    "2fa",
		"iat":     time.Now().Unix(),
		"scope":   "all",
		"token":   "access",
	}
	twoFATimeout := time.Duration(config.GetConfig().GetInt64("auth.key2fa_token_timeout") * int64(time.Second))
	token, err := middleware.GenerateJWT(jwtSecret, claims, time.Now().Add(time.Duration(twoFATimeout)))
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
		"mobile":  userRef.Mobile.String,
		"country": userRef.CountryCode.String,
		"type":    "2fa",
		"iat":     time.Now().Unix(),
		"scope":   "all",
		"token":   "refresh",
	}
	twoFARefreshTimeout := time.Duration(config.GetConfig().GetInt64("auth.key2fa_refresh_timeout") * int64(time.Second))
	refreshToken, err := middleware.GenerateJWT(jwtSecret, refreshClaims, time.Now().Add(twoFARefreshTimeout))
	if err != nil {
		logger.Log(c).Error("failed to generate JWT token", zap.Error(err))
		response.Error = append(response.Error, "Failed to generate JWT token")
		response.Description = "Internal server error"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	response.Description = "mpin login successfull"
	response.Data = &dto.VerifyMPINData{
		AccessToken:        token,
		RefreshToken:       refreshToken,
		AccessTokenExpiry:  refreshClaims["iat"].(int64) + 300,   // 5 minutes
		RefreshTokenExpiry: refreshClaims["iat"].(int64) + 86400, // 24 hours
	}

	c.JSON(http.StatusOK, response)
}
