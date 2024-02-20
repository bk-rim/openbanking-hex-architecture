package controller

import (
	"net/http"

	"github.com/bk-rim/openbanking/domain/service"
	"github.com/bk-rim/openbanking/model"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService service.PaymentService
}

func NewPaymentController(paymentService *service.PaymentService) *PaymentController {
	return &PaymentController{paymentService: *paymentService}
}

func (pc *PaymentController) GetAllPayments(c *gin.Context) {
	payments, err := pc.paymentService.GetAllPayments()
	if err != nil {
		c.IndentedJSON(err.Status, gin.H{"error": err.Message})
		return
	}
	c.IndentedJSON(http.StatusOK, payments)
}

func (pc *PaymentController) GetPaymentsByIbanCdt(c *gin.Context) {
	iban := c.Param("iban")
	payments, err := pc.paymentService.GetPaymentsByIban(iban, "creditor")
	if err != nil {
		c.IndentedJSON(err.Status, gin.H{"error": err.Message})
		return
	}
	c.IndentedJSON(http.StatusOK, payments)
}

func (pc *PaymentController) GetPaymentsByIbanDbt(c *gin.Context) {
	iban := c.Param("iban")
	payments, err := pc.paymentService.GetPaymentsByIban(iban, "debtor")
	if err != nil {
		c.IndentedJSON(err.Status, gin.H{"error": err.Message})
		return
	}
	c.IndentedJSON(http.StatusOK, payments)
}

func (pc *PaymentController) PaymentHandler(responseChannel chan<- model.PaymentResponse) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payment model.Payment
		if err := c.ShouldBindJSON(&payment); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := pc.paymentService.HandlePayment(payment, responseChannel)
		if err != nil {
			c.IndentedJSON(err.Status, gin.H{"error": err.Message})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "payment is being processed"})

	}
}
