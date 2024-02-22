package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/bk-rim/openbanking/model"
)

type BankService struct {
	bankRepository    IBankRepository
	fileCsvRepository IFileCsvRepository
}

type IBankRepository interface {
	UpdatePaymentStatus(response model.PaymentResponse) error
}

type IFileCsvRepository interface {
	Save(response model.PaymentResponse)
}

func NewBankService(bankRepository IBankRepository, fileCsvRepository IFileCsvRepository) *BankService {
	return &BankService{bankRepository: bankRepository, fileCsvRepository: fileCsvRepository}
}

func (bs *BankService) notifyClient(paymentResponse model.PaymentResponse, webHookUrl string) int {

	payload, err := json.Marshal(paymentResponse)
	username := os.Getenv("USER_NAME")
	password := os.Getenv("PASSWORD")
	if err != nil {
		fmt.Println("error marshalling payment response")
		return 0
	}

	req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println("error creating request")
		return 0
	}

	client := &http.Client{}

	req.Header.Set("Basic-Auth", username+":"+password)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error sending request", err.Error())
		return 0
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		fmt.Println("webhook sent successfully")
	} else {
		fmt.Println("error sending webhook")
	}

	return res.StatusCode
}

func (bs *BankService) HandleBankResponses(responseChannel <-chan model.PaymentResponse, webhookUrl string) {

	for response := range responseChannel {

		if err := bs.bankRepository.UpdatePaymentStatus(response); err != nil {
			fmt.Println("error updating status in database")
		}
		bs.fileCsvRepository.Save(response)
		if response.Status != "PENDING" {
			bankService := BankService{}
			bankService.notifyClient(response, webhookUrl)
		}
		fmt.Printf("Payment ID: %s, Status: %s\n", response.Id, response.Status)
	}
}
