package server

import (
	"fmt"
	"net/http"
	"time"

	// auth "bitbucket.org/liquide-life/be-auth-go"

	"github.com/gin-gonic/gin"

	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/constant"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"github.com/sauravkuila/mergemoney_assessment/pkg/middleware"
	"github.com/sauravkuila/mergemoney_assessment/pkg/service"
	"go.uber.org/zap"
)

// starts the server router in a go routine.
//
//	uses a gin Engine as handler function to a native http router
//	tracks the router instance globally
//		srv = &http.Server{
//			Addr:    <config ip>:<config port>,
//		 	Handler: &gin.Emgine{},
//		}
func startRouter(obj service.ServiceItf) {
	srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GetConfig().GetInt("server.http-port")),
		Handler: getRouter(obj), //getRouter set the api specs for version-1 routes
	}
	// run api router
	logger.Log().Info("starting router")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log().Fatal("Error starting server", zap.Error(err))
		}
	}()
}

func getRouter(serviceObj service.ServiceItf) *gin.Engine {
	if config.GetConfig().GetString("env") != string(constant.DEVELOPMENT) {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	router.Use(traceLogger(logger.Log()))
	router.Use(gin.Recovery())

	//health check
	router.GET("/health", serviceObj.Health)

	// serve sample UI pages (static) -- use parent path because the server runs from cmd/api
	router.Static("/ui", "../../sample_UI")
	// root -> UI index
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/ui/index.html")
	})

	loginGroup := router.Group("/v1")
	{
		loginGroup.GET("/generateOTP", serviceObj.GetV1Object().GenerateOTP)
		loginGroup.POST("/verifyOTP", serviceObj.GetV1Object().VerifyOTP)
		loginGroup.POST("/resetMPIN", serviceObj.GetV1Object().ResetMPIN)

		oneFAGroup := loginGroup.Group("")
		{
			oneFAGroup.Use(middleware.OneFAMiddleware())
			oneFAGroup.POST("/setMPIN", serviceObj.GetV1Object().SetMPIN)
			oneFAGroup.POST("/1fa/refresh", serviceObj.GetV1Object().Refresh1FA)
			oneFAGroup.POST("/verifyMPIN", serviceObj.GetV1Object().VerifyMPIN)
			oneFAGroup.GET("/accounts", serviceObj.GetV1Object().GetAccounts) // get accounts against mobile number
		}

		twoFAGroup := loginGroup.Group("")
		{
			twoFAGroup.Use(middleware.TwoFAMiddleware())
			twoFAGroup.POST("/2fa/refresh", serviceObj.GetV1Object().Refresh2FA)
			transferGroup := twoFAGroup.Group("transfer")
			{
				transferGroup.POST("", serviceObj.GetV1Object().Transfer)
				transferGroup.POST("/confirm", serviceObj.GetV1Object().TransferConfirm)
			}
		}
	}

	return router
}

func traceLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		uID := c.GetHeader("x-request-id")
		if uID == "" {
			uID = c.GetString(config.REQUESTID)
		}
		//generate a requestid if cintext does not have one
		if uID == "" {
			uID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Set(config.REQUESTID, uID)

		c.Next()

		if c.FullPath() != "/health" {
			latency := time.Since(start).Milliseconds()
			// userID := c.GetString(config.USERID)
			uID := c.GetString(config.REQUESTID)
			logger.Info("request_response",
				zap.String("path", path),
				zap.String("requestID", uID),
				// zap.String("userId", userID),
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Int64("latency", latency),
			)
		}
	}
}
