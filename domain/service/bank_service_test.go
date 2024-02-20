package service

import (
	"sync"
	"testing"

	"fmt"
	"net/http"

	"github.com/bk-rim/openbanking/model"
	"github.com/stretchr/testify/assert"
)

type MockBankRepository struct{}

func (*MockBankRepository) UpdatePaymentStatus(response model.PaymentResponse) error {
	return nil
}

func serverWebhook() {
	http.HandleFunc("/webhook", handleWebhook)
	http.ListenAndServe(":8080", nil)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "Webhook received")

}
func TestBankService_notifyClient(t *testing.T) {
	mockBankRepository := &MockBankRepository{}
	bankService := NewBankService(mockBankRepository)
	paymentResponse := model.PaymentResponse{Id: "JBXXXZZ", Status: "PROCESSED"}
	go serverWebhook()
	statusCode := bankService.notifyClient(paymentResponse, "http://localhost:8080/webhook")

	assert.Equal(t, 200, statusCode)
}
func TestBankService_HandleBankResponses(t *testing.T) {
	responseChannel := make(chan model.PaymentResponse)

	webhookURL := "http://localhost:8080/webhook"

	mockBankRepository := &MockBankRepository{}
	bankService := NewBankService(mockBankRepository)

	var wg sync.WaitGroup
	wg.Add(1)
	paymentResponse := model.PaymentResponse{Id: "JBXXXZZ", Status: "PROCESSED"}
	go func() {
		defer wg.Done()
		bankService.HandleBankResponses(responseChannel, webhookURL)

	}()

	responseChannel <- paymentResponse

	close(responseChannel)
	wg.Wait()

	for response := range responseChannel {
		assert.Equal(t, paymentResponse, response)
	}

}
