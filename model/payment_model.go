package model

import (
	"github.com/bk-rim/openbanking/helper"
)

type Payment struct {
	DebtorIBAN           string  `json:"debtor_iban" validate:"required"`
	DebtorName           string  `json:"debtor_name" validate:"required,min=3,max=30"`
	CreditorIBAN         string  `json:"creditor_iban" validate:"required"`
	CreditorName         string  `json:"creditor_name" validate:"required,min=3,max=30"`
	Amount               float64 `json:"ammount" validate:"required,number"`
	IdempotencyUniqueKey string  `json:"idempotency_unique_key" validate:"required"`
	Status               string  `json:"status"`
}

type PaymentResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func (p *Payment) ValidateIBAN() bool {
	return helper.IsValidIban(p.DebtorIBAN) && helper.IsValidIban(p.CreditorIBAN)
}
