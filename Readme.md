
## Description and Approach 
```This is an example of openBanking api that allows you to issue bank payments.```

## Stack
- Golang
- Gin
- Sqlite
- testify

## Architecture & Project structure
- hexagonal architecture

```md
.
├── api
│   ├── controller
│   │   ├── payment_controller.go
│   │   ├── webhook_controller.go
│   │   └── webhook_controller_test.go
│   └── middleware
│       └── auth.go
├── domain
│   └── service
│       ├── bank_service.go
│       ├── bank_service_test.go
│       ├── generatePaymentXml.go
│       ├── payment_service.go
│       ├── payment_service_test.go
│       └── simulateBankProcessing.go
├── go.mod
├── go.sum
├── helper
│   └── isValidIban.go
├── main.go
├── model
│   └── payment_model.go
├── Readme.md
└── repository
    ├── file
    │   ├── save_csv.go
    │   └── save_xml.go
    ├── sqlite
    │   ├── bank_repository.go
    │   ├── init_DB.go
    │   └── payment_repository.go
    └── utils
        └── IKey_repository.go

```

## Features

- generate bank payment
- Simulate bank processing
- handle bank response with a webhook
- webhookclient endpoint
- generate xml for payment 
- generate csv for payment with payment status
- get all payment 
- get payment filter by iban debitor
- get payment filter by iban creditor

## endpoints

- /bank
    - /payment (send payment)
    - /payments/ibanCdt/:iban (get payments by iban creditor)
    - /payments/ibanDbt/:iban (get payments by iban debitor)
    - /allpayments (get all payments)

- /client
    - /webhook (handle webhook)