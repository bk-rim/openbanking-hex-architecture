package service

import (
	"github.com/bk-rim/openbanking/helper"
	"github.com/bk-rim/openbanking/model"
	"github.com/go-playground/validator/v10"
)

type PaymentService struct {
	paymentRepository IPaymentRepository
	iKeyRepository    IdempotentKey
	fileXmlRepository IFileXmlRepository
}

type IPaymentRepository interface {
	FindAll() ([]model.Payment, error)
	FindByIbanCdt(iban string) ([]model.Payment, error)
	FindByIbanDbt(iban string) ([]model.Payment, error)
	Save(payment model.Payment) error
}

type IdempotentKey interface {
	Generate() (string, error)
}

type IFileXmlRepository interface {
	Save(xmlData []byte, idempotencyKey string, responseChannel chan<- model.PaymentResponse)
}

type ErrorService struct {
	Message string
	Status  int
}

func NewPaymentService(paymentRepository IPaymentRepository, fileXmlRepository IFileXmlRepository, generateIdempotentKey IdempotentKey) *PaymentService {
	return &PaymentService{paymentRepository: paymentRepository, fileXmlRepository: fileXmlRepository, iKeyRepository: generateIdempotentKey}
}

func (ps *PaymentService) GetAllPayments() ([]model.Payment, *ErrorService) {

	payments, err := ps.paymentRepository.FindAll()
	if err != nil {
		return nil, &ErrorService{err.Error(), 500}
	}
	return payments, nil
}

func (ps *PaymentService) GetPaymentsByIban(iban string, actor string) ([]model.Payment, *ErrorService) {

	if !helper.IsValidIban(iban) {
		return nil, &ErrorService{"iban not valid", 400}
	}

	switch actor {
	case "creditor":
		payments, err := ps.paymentRepository.FindByIbanCdt(iban)
		if err != nil {
			return nil, &ErrorService{err.Error(), 500}
		}
		return payments, nil
	case "debtor":
		payments, err := ps.paymentRepository.FindByIbanDbt(iban)
		if err != nil {
			return nil, &ErrorService{err.Error(), 500}
		}
		return payments, nil
	default:
		return nil, &ErrorService{"actor not valid", 400}
	}

}

func (ps *PaymentService) HandlePayment(payment model.Payment, responseChannel chan<- model.PaymentResponse) *ErrorService {

	if err := payment.ValidateIBAN(); !err {
		return &ErrorService{"iban not valid", 400}
	}

	iKey, err := ps.iKeyRepository.Generate()
	if err != nil {
		return &ErrorService{err.Error(), 500}
	}

	payment.IdempotencyUniqueKey = iKey
	validate := validator.New()
	if err := validate.Struct(payment); err != nil {
		return &ErrorService{err.Error(), 400}
	}

	xmlData, err := generateXML(payment)
	if err != nil {
		return &ErrorService{err.Error(), 500}
	}

	if err := ps.paymentRepository.Save(payment); err != nil {
		return &ErrorService{err.Error(), 500}
	}

	go ps.fileXmlRepository.Save(xmlData, payment.IdempotencyUniqueKey, responseChannel)

	go simulateBankProcessing(payment.IdempotencyUniqueKey, responseChannel)

	return nil
}
