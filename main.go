package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/bk-rim/openbanking/api/middleware"
	"github.com/bk-rim/openbanking/domain/service"
	"github.com/bk-rim/openbanking/repository/file"
	database "github.com/bk-rim/openbanking/repository/sqlite"
	"github.com/bk-rim/openbanking/repository/utils"

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
	database.InitSqlite()

	responseChannel := make(chan model.PaymentResponse)
	bankRepository := &database.BankRepository{}
	iKeyRepository := &utils.IKeyRepository{}
	fileXmlRepository := &file.FileXmlRepository{}
	fileCsvRepository := &file.FileCsvRepository{}
	bankService := service.NewBankService(bankRepository, fileCsvRepository)
	paymentRepository := &database.PaymentRepository{}
	paymentService := service.NewPaymentService(paymentRepository, fileXmlRepository, iKeyRepository)
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
