package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bk-rim/openbanking/api/controller"
	"github.com/bk-rim/openbanking/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWebhookHandler(t *testing.T) {
	r := gin.Default()
	r.POST("/webhook", controller.WebhookHandler)
	paymentResponse := model.PaymentResponse{
		Id:     "123",
		Status: "success",
	}

	payload, err := json.Marshal(paymentResponse)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWebHookHandlerWithBadRequest(t *testing.T) {
	r := gin.Default()
	r.POST("/webhook", controller.WebhookHandler)
	paymentResponse := "bad request"

	payload, err := json.Marshal(paymentResponse)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
