package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/bk-rim/openbanking/api/middleware"
	"github.com/bk-rim/openbanking/domain/service"
	repository "github.com/bk-rim/openbanking/repository/sqlite"

	"github.com/bk-rim/openbanking/model"

	"github.com/bk-rim/openbanking/api/controller"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

func main() {

	loadEnv()
	repository.InitSqlite()

	responseChannel := make(chan model.PaymentResponse)
	bankRepository := &repository.SqliteBankRepository{}
	bankService := service.NewBankService(bankRepository)
	paymentRepository := &repository.SqlitePaymentRepository{}
	paymentService := service.NewPaymentService(paymentRepository)
	paymentController := controller.NewPaymentController(paymentService)

	go bankService.HandleBankResponses(responseChannel, "http://localhost:8080/client/webhook")

	r := gin.Default()
	bankRoute := r.Group("/bank")
	bankRoute.Use(middleware.AuthMiddleware)
	bankRoute.POST("/payment", paymentController.PaymentHandler(responseChannel))
	bankRoute.GET("/allpayments", paymentController.GetAllPayments)
	bankRoute.GET("/payments/ibanCdt/:iban", paymentController.GetPaymentsByIbanCdt)
	bankRoute.GET("/payments/ibanDbt/:iban", paymentController.GetPaymentsByIbanDbt)

	clientRoute := r.Group("/client")
	clientRoute.POST("/webhook", controller.WebhookHandler)

	r.Run(":8080")
}
