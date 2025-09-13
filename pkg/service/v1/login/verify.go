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

func (obj *loginSt) Refresh1FA(c *gin.Context) {
	var (
		request  dto.Refresh1FARequest
		response dto.Refresh1FAAResponse
	)

	//bind the request
	if err := c.BindJSON(&request); err != nil {
		logger.Log(c).Error("failed to bind request", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//fetch the claims from context
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, &gin.H{"error": "unauthorized"})
		return
	}
	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "invalid claims type"})
		return
	}

	jwtSecret := config.GetConfig().GetString("auth.key")

	//validate the request sent refresh token
	_, err := middleware.ValidateJWT(request.RefreshToken, jwtSecret)
	if err != nil {
		logger.Log(c).Error("failed to validate refresh token", zap.Error(err))
		response.Error = append(response.Error, "Invalid refresh token")
		response.Description = "Invalid refresh token"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//generate an access token
	accessClaims := map[string]interface{}{
		"user_id": claimsMap["user_id"],
		"mobile":  claimsMap["mobile"],
		"country": claimsMap["country"],
		"type":    claimsMap["type"],
		"iat":     time.Now().Unix(),
		"scope":   claimsMap["scope"],
		"token":   claimsMap["token"],
	}
	oneFATokenTimeout := time.Duration(config.GetConfig().GetInt64("auth.key1fa_token_timeout") * int64(time.Second))
	token, err := middleware.GenerateJWT(jwtSecret, accessClaims, time.Now().Add(oneFATokenTimeout))
	if err != nil {
		logger.Log(c).Error("failed to generate JWT token", zap.Error(err))
		response.Error = append(response.Error, "Failed to generate JWT token")
		response.Description = "Internal server error"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	response.Description = "access token refreshed successfully"
	response.Data = &dto.Refresh1FA{
		OneFAAccessToken: token,
	}

	c.JSON(http.StatusOK, response)
}

func (obj *loginSt) Refresh2FA(c *gin.Context) {
	var (
		request  dto.Refresh2FARequest
		response dto.Refresh2FAResponse
	)

	//bind the request
	if err := c.BindJSON(&request); err != nil {
		logger.Log(c).Error("failed to bind request", zap.Error(err))
		response.Error = append(response.Error, err.Error())
		response.Description = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//fetch the claims from context
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, &gin.H{"error": "unauthorized"})
		return
	}
	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "invalid claims type"})
		return
	}

	jwtSecret := config.GetConfig().GetString("auth.key2fa")

	//validate the request sent refresh token
	_, err := middleware.ValidateJWT(request.RefreshToken, jwtSecret)
	if err != nil {
		logger.Log(c).Error("failed to validate refresh token", zap.Error(err))
		response.Error = append(response.Error, "Invalid refresh token")
		response.Description = "Invalid refresh token"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//generate an access token
	accessClaims := map[string]interface{}{
		"user_id": claimsMap["user_id"],
		"mobile":  claimsMap["mobile"],
		"country": claimsMap["country"],
		"type":    claimsMap["type"],
		"iat":     time.Now().Unix(),
		"scope":   claimsMap["scope"],
		"token":   claimsMap["token"],
	}
	twoFATimeout := time.Duration(config.GetConfig().GetInt64("auth.key2fa_token_timeout") * int64(time.Second))
	token, err := middleware.GenerateJWT(jwtSecret, accessClaims, time.Now().Add(twoFATimeout))
	if err != nil {
		logger.Log(c).Error("failed to generate JWT token", zap.Error(err))
		response.Error = append(response.Error, "Failed to generate JWT token")
		response.Description = "Internal server error"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	response.Description = "access token refreshed successfully"
	response.Data = &dto.Refresh2FA{
		AccessToken: token,
	}

	c.JSON(http.StatusOK, response)
}
