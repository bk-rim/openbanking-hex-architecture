package file

import (
	"os"

	"fmt"

	"path/filepath"

	"github.com/bk-rim/openbanking/model"
)

type FileXmlRepository struct{}

func (fr *FileXmlRepository) Save(xmlData []byte, idempotencyKey string, responseChannel chan<- model.PaymentResponse) {

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
