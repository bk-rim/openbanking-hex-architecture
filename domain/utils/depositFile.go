package utils

import (
	"os"

	"fmt"

	"path/filepath"

	"encoding/csv"

	"github.com/bk-rim/openbanking/model"
)

func DepositPaymentInBankFolder(xmlData []byte, idempotencyKey string, responseChannel chan<- model.PaymentResponse) {

	bankFolderPath := os.Getenv("BANK_FOLDER_PATH")
	if bankFolderPath == "" {
		fmt.Println("BANK_FOLDER_PATH environment variable is not set")
		os.Exit(1)
	}

	if err := os.MkdirAll(bankFolderPath, os.ModePerm); err != nil {
		responseChannel <- model.PaymentResponse{Id: idempotencyKey, Status: "ERROR", Error: err.Error()}
		return
	}

	xmlFilePath := filepath.Join(bankFolderPath, fmt.Sprintf("%s.xml", "payment-"+idempotencyKey))

	if err := os.WriteFile(xmlFilePath, xmlData, os.ModePerm); err != nil {
		responseChannel <- model.PaymentResponse{Id: idempotencyKey, Status: "ERROR", Error: err.Error()}
		return
	}

	responseChannel <- model.PaymentResponse{Id: idempotencyKey, Status: "PENDING", Error: ""}

}

func DepositResponseInBankFolder(response model.PaymentResponse) {
	bankFolderPath := os.Getenv("BANK_FOLDER_PATH")
	if bankFolderPath == "" {
		fmt.Println("BANK_FOLDER_PATH environment variable is not set")
		os.Exit(1)
	}

	csvData := [][]string{
		{"id", "status"},
		{response.Id, response.Status},
	}

	if err := os.MkdirAll(bankFolderPath, os.ModePerm); err != nil {
		fmt.Printf("Error while creating bank folder: %s\n", err.Error())
		return
	}

	csvFilePath := filepath.Join(bankFolderPath, fmt.Sprintf("%s.csv", "bankResponse-"+response.Id))

	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		fmt.Printf("Error while creating CSV file: %s\n", err.Error())
		return
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	if err := csvWriter.WriteAll(csvData); err != nil {
		fmt.Printf("Error while writing CSV file: %s\n", err.Error())
		return
	}

}
