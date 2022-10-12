package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
)

func WalletBalance(c *gin.Context) {

	userID := c.GetUint("id")
	var balance int
	database.DB.Raw("SELECT balance FROM wallets WHERE users_id = ?", userID).Scan(&balance)
	c.JSON(200, gin.H{
		"status":  true,
		"Balance": balance,
		"UserID":  userID,
	})
}

func WalletInfo(c *gin.Context) {

	UsersID := c.GetUint("id")
	var history []models.Wallethistory
	database.DB.Where("wallethistories.users_id = ?", UsersID).Find(&history)

	for _, i := range history {
		c.JSON(200, gin.H{
			"status": true,
			"Date":   i.CreatedAt,
			"Debit":  i.Debit,
			"Credit": i.Credit,
			"UserID": i.UsersID,
		})
	}

}
