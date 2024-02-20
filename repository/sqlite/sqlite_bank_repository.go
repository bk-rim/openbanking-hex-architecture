package repository

import (
	"github.com/bk-rim/openbanking/model"
)

type SqliteBankRepository struct{}

func (br *SqliteBankRepository) UpdatePaymentStatus(response model.PaymentResponse) error {
	_, err := db.Exec("UPDATE payments SET status = ? WHERE idempotency = ?", response.Status, response.Id)
	if err != nil {
		return err
	}
	return nil
}
