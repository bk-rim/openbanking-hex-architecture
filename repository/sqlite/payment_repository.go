package database

import (
	"github.com/bk-rim/openbanking/model"
)

type PaymentRepository struct{}

func (pr *PaymentRepository) Save(payment model.Payment) error {
	_, err := db.Exec("INSERT INTO payments (idempotency, debtor_iban, debtor_name, creditor_iban, creditor_name, amount, status) VALUES (?, ?, ?, ?, ?, ?, ?)", payment.IdempotencyUniqueKey, payment.DebtorIBAN, payment.DebtorName, payment.CreditorIBAN, payment.CreditorName, payment.Amount, "PENDING")
	if err != nil {
		return err
	}
	return nil
}

func (pr *PaymentRepository) FindAll() ([]model.Payment, error) {
	rows, err := db.Query("SELECT idempotency, debtor_iban, debtor_name, creditor_iban, creditor_name, amount, status FROM payments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var payment model.Payment
		if err := rows.Scan(&payment.IdempotencyUniqueKey,
			&payment.DebtorIBAN,
			&payment.DebtorName,
			&payment.CreditorIBAN,
			&payment.CreditorName,
			&payment.Amount,
			&payment.Status); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (pr *PaymentRepository) FindByIbanCdt(iban string) ([]model.Payment, error) {
	rows, err := db.Query("SELECT idempotency, debtor_iban, debtor_name, creditor_iban, creditor_name, amount, status FROM payments WHERE creditor_iban = $1", iban)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var payment model.Payment
		if err := rows.Scan(&payment.IdempotencyUniqueKey,
			&payment.DebtorIBAN,
			&payment.DebtorName,
			&payment.CreditorIBAN,
			&payment.CreditorName,
			&payment.Amount,
			&payment.Status); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (pr *PaymentRepository) FindByIbanDbt(iban string) ([]model.Payment, error) {
	rows, err := db.Query("SELECT idempotency, debtor_iban, debtor_name, creditor_iban, creditor_name, amount, status FROM payments WHERE debtor_iban = $1", iban)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var payment model.Payment
		if err := rows.Scan(&payment.IdempotencyUniqueKey,
			&payment.DebtorIBAN,
			&payment.DebtorName,
			&payment.CreditorIBAN,
			&payment.CreditorName,
			&payment.Amount,
			&payment.Status); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}
