package controller

import (
	"net/http"

	"github.com/bk-rim/openbanking/model"

	"github.com/gin-gonic/gin"
)

func WebhookHandler(c *gin.Context) {
	var paymentResponse model.PaymentResponse

	if err := c.ShouldBindJSON(&paymentResponse); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "webhook received"})
}
