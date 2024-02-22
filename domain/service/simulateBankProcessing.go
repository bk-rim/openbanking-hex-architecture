package service

import (
	"time"

	"github.com/bk-rim/openbanking/model"
)

func simulateBankProcessing(idempotencyKey string, responseChannel chan<- model.PaymentResponse) {
	time.Sleep(5 * time.Second)

	status := "PROCESSED"
	responseChannel <- model.PaymentResponse{Id: idempotencyKey, Status: status}
}
