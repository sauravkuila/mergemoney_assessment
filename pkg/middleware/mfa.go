package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/constant"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"go.uber.org/zap"
)

func OneFAMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response dto.AuthResponse
		response = dto.AuthResponse{
			Status:  false,
			Errors:  []string{"UnAuthorized: invalid auth token"},
			Message: "Authentication failed",
		}

		// Access Token will be appended with Bearer, need to get only the token
		accessTokenString := c.Request.Header.Get(constant.AUTH_ACCESS_TOKEN)
		if accessTokenString == "" {
			logger.Log().Error("auth token incorrect")
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		claims, err := ValidateJWT(accessTokenString, config.GetConfig().GetString("auth.key"))
		if err != nil && err.Error() != "token expired" {
			logger.Log().Error("auth token validation failed", zap.Error(err))
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		} else if err != nil && err.Error() == "token expired" {

			logger.Log(c).Info("claims are available but token is expired", zap.Error(err), zap.Any("claims", claims))
			//check if token claim is a refresh or access token
			//return failed if refresh token is expired
			if v, f := claims["token"]; f && v == "refresh" {
				logger.Log().Error("refresh token expired", zap.Error(err))
				response.Message = "refresh token expired"
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			}

			// if token is access and is expired, only allow for 1fa/refresh path
			if c.Request.URL.Path != "/v1/1fa/refresh" {
				logger.Log().Error("access token expired", zap.Error(err))
				response.Message = "access token expired"
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			}
		}

		c.Set(config.AUTHORIZATION, accessTokenString)
		c.Set("claims", claims)
		c.Set(config.USERID, claims["user_id"])
		c.Set(config.MOBILE, claims["mobile"])
		c.Next()
	}
}

func TwoFAMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response dto.AuthResponse
		response = dto.AuthResponse{
			Status:  false,
			Errors:  []string{"UnAuthorized: invalid auth token"},
			Message: "Authentication failed",
		}

		// Access Token will be appended with Bearer, need to get only the token
		accessTokenString := c.Request.Header.Get(constant.AUTH_ACCESS_TOKEN)
		if accessTokenString == "" {
			logger.Log().Error("auth token incorrect")
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		claims, err := ValidateJWT(accessTokenString, config.GetConfig().GetString("auth.key2fa"))
		if err != nil && err.Error() != "token expired" {
			logger.Log().Error("auth token validation failed", zap.Error(err))
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		} else if err != nil && err.Error() == "token expired" {

			logger.Log(c).Info("claims are available but token is expired", zap.Error(err), zap.Any("claims", claims))
			//check if token claim is a refresh or access token
			//return failed if refresh token is expired
			if v, f := claims["token"]; f && v == "refresh" {
				logger.Log().Error("refresh token expired", zap.Error(err))
				response.Message = "refresh token expired"
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			}

			// if token is access and is expired, only allow for 1fa/refresh path
			if c.Request.URL.Path != "/v1/2fa/refresh" {
				logger.Log().Error("access token expired", zap.Error(err))
				response.Message = "access token expired"
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			}
		}

		c.Set(config.AUTHORIZATION, accessTokenString)
		c.Set("claims", claims)
		c.Set(config.USERID, claims["user_id"])
		c.Set(config.MOBILE, claims["mobile"])
		c.Next()
	}
}
