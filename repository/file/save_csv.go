package file

import (
	"os"

	"fmt"

	"path/filepath"

	"encoding/csv"

	"github.com/bk-rim/openbanking/model"
)

type FileCsvRepository struct{}

func (fr *FileCsvRepository) Save(response model.PaymentResponse) {
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
