package service_test

import (
	"log"
	"testing"

	"github.com/bk-rim/openbanking/domain/service"
	"github.com/bk-rim/openbanking/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type MockPaymentRepository struct{}

type MockIKeyRepository struct{}

type MockFileXmlRepository struct{}

func (*MockFileXmlRepository) Save(xmlData []byte, idempotencyKey string, responseChannel chan<- model.PaymentResponse) {
	return
}

func (*MockIKeyRepository) Generate() (string, error) {
	return "JXJ" + "XXXZ", nil
}

func (*MockPaymentRepository) FindAll() ([]model.Payment, error) {
	return []model.Payment{
		{DebtorIBAN: "AE1234567890", CreditorIBAN: "BE0987654321", Amount: 100.0},
		{DebtorIBAN: "CE1234567890", CreditorIBAN: "DE0987654321", Amount: 200.0},
	}, nil
}

func (*MockPaymentRepository) FindByIbanCdt(iban string) ([]model.Payment, error) {
	allPayments := []model.Payment{
		{DebtorIBAN: "AE1234567890", CreditorIBAN: "BE0987654321", Amount: 100.0},
		{DebtorIBAN: "CE1234567890", CreditorIBAN: "DE0987654321", Amount: 200.0},
	}

	var payments []model.Payment
	for _, payment := range allPayments {
		if payment.CreditorIBAN == iban {
			payments = append(payments, payment)
			return payments, nil
		}
	}
	return []model.Payment{}, nil
}

func (*MockPaymentRepository) FindByIbanDbt(iban string) ([]model.Payment, error) {
	allPayments := []model.Payment{
		{DebtorIBAN: "AE1234567890", CreditorIBAN: "BE0987654321", Amount: 100.0},
		{DebtorIBAN: "CE1234567890", CreditorIBAN: "DE0987654321", Amount: 200.0},
	}

	var payments []model.Payment
	for _, payment := range allPayments {
		if payment.DebtorIBAN == iban {
			payments = append(payments, payment)
			return payments, nil
		}
	}
	return []model.Payment{}, nil
}

func (*MockPaymentRepository) Save(payment model.Payment) error {
	return nil
}

func TestPaymentService_GetAllPayments(t *testing.T) {

	mockPaymentRepository := &MockPaymentRepository{}

	mockFileXmlRepository := &MockFileXmlRepository{}

	mockIKeyRepository := &MockIKeyRepository{}

	paymentService := service.NewPaymentService(mockPaymentRepository, mockFileXmlRepository, mockIKeyRepository)

	expectedPayments := []model.Payment{
		{DebtorIBAN: "AE1234567890", CreditorIBAN: "BE0987654321", Amount: 100.0},
		{DebtorIBAN: "CE1234567890", CreditorIBAN: "DE0987654321", Amount: 200.0},
	}

	payments, err := paymentService.GetAllPayments()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	assert.Equal(t, expectedPayments, payments)

}

func TestPaymentService_GetPaymentsByIban(t *testing.T) {

	mockPaymentRepository := &MockPaymentRepository{}

	mockFileXmlRepository := &MockFileXmlRepository{}

	mockIKeyRepository := &MockIKeyRepository{}

	paymentService := service.NewPaymentService(mockPaymentRepository, mockFileXmlRepository, mockIKeyRepository)

	iban_1 := "DE0987654321"
	actor_1 := "creditor"

	expectedPayments_1 := []model.Payment{
		{DebtorIBAN: "CE1234567890", CreditorIBAN: "DE0987654321", Amount: 200.0},
	}

	iban_2 := "AE1234567890"
	actor_2 := "debtor"

	expectedPayments_2 := []model.Payment{
		{DebtorIBAN: "AE1234567890", CreditorIBAN: "BE0987654321", Amount: 100.0},
	}

	payments_1, err := paymentService.GetPaymentsByIban(iban_1, actor_1)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	payments_2, err := paymentService.GetPaymentsByIban(iban_2, actor_2)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	assert.Equal(t, expectedPayments_1, payments_1)
	assert.Equal(t, expectedPayments_2, payments_2)

}

func TestPaymentService_HandlePayment(t *testing.T) {

	errLoad := godotenv.Load("../../.env")
	if errLoad != nil {
		log.Fatal("Error loading .env file", errLoad)
	}
	mockPaymentRepository := &MockPaymentRepository{}

	mockFileXmlRepository := &MockFileXmlRepository{}

	mockIKeyRepository := &MockIKeyRepository{}

	paymentService := service.NewPaymentService(mockPaymentRepository, mockFileXmlRepository, mockIKeyRepository)

	payment := model.Payment{
		DebtorIBAN:   "AE1234567890",
		DebtorName:   "John Doe",
		CreditorIBAN: "BE0987654321",
		CreditorName: "Jane Doe",
		Amount:       100.0,
	}

	responseChannel := make(chan model.PaymentResponse)

	err := paymentService.HandlePayment(payment, responseChannel)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	} else {
		paymentResponse := <-responseChannel
		assert.Equal(t, "PROCESSED", paymentResponse.Status)
	}

}
