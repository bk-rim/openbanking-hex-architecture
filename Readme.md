
## Description and Approach 
```This is an example of openBanking api that allows you to issue bank payments.```

## Stack
- Golang
- Gin
- Sqlite
- testify

## Architecture
- hexagonal architecture

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