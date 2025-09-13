package account

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/external/accountaggregator"
	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"go.uber.org/zap"
)

func (obj *accountSt) GetAccounts(c *gin.Context) {
	var (
		response dto.GetAccountsResponse
	)
	//get mobile number from context
	mobile := c.GetString(config.MOBILE)
	userId := c.GetString(config.USERID)

	//fetch user accounts from DB, ask AA if DB accounts are not found
	userAccountsInDb, err := obj.DB.GetUserAccountsByUserId(c, userId)
	if err == nil {
		logger.Log(c).Info("fetched user accounts from DB", zap.String("mobile", mobile), zap.Int("count", len(userAccountsInDb)))
		//check if accounts are found and not old
		if len(userAccountsInDb) > 0 {
			response.Data = make([]dto.UserAccount, 0)
			for _, acc := range userAccountsInDb {
				data := acc.ToAggregatorAccount()
				response.Data = append(response.Data, data)
			}
			response.Status = true
			response.Description = "Accounts fetched successfully from DB"
			c.JSON(http.StatusOK, response)
			return
		}
	}

	//call AA api to get accounts
	data, err := accountaggregator.GetAccountsAgainstMobile(c, mobile, obj.utilsObj)
	if err != nil {
		logger.Log(c).Error("error in fetching accounts from AA", zap.Error(err))
		response.Error = append(response.Error, fmt.Sprintf("error in fetching accounts from AA: %v", err))
		response.Description = "Failed to fetch accounts"
		c.JSON(http.StatusFailedDependency, response)
		return
	}

	//check if no accounts found
	if len(data) == 0 {
		logger.Log(c).Info("no accounts found for user", zap.String("mobile", mobile))
		response.Description = "No accounts found"
		c.JSON(http.StatusNoContent, response)
		return
	}

	// save accounts to DB for future reference if needed
	if dbSavedAccounts, err := obj.DB.SaveUserAccounts(c, userId, data); err != nil {
		logger.Log(c).Error("error in saving accounts to DB", zap.Error(err))
		response.Error = append(response.Error, fmt.Sprintf("error in saving accounts to DB: %v", err))
		response.Description = "Failed to save accounts"
		c.JSON(http.StatusInternalServerError, response)
		return
	} else {
		response.Data = make([]dto.UserAccount, 0)
		sids := make([]int64, 0)
		for _, acc := range dbSavedAccounts {
			data := acc.ToAggregatorAccount()
			response.Data = append(response.Data, data)
			sids = append(sids, acc.Sid.Int64)
		}
		logger.Log(c).Info("saved user accounts to DB", zap.String("mobile", mobile), zap.Int64s("serial_ids", sids))
	}

	// return accounts
	// response.Data = append(response.Data, data...)
	response.Status = true
	response.Description = "Accounts fetched successfully"

	c.JSON(http.StatusOK, response)
}
