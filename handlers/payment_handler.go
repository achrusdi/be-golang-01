package handlers

import (
	"encoding/json"
	"log"
	"mnc-be-golang/usecases"
	"net/http"
)

type PaymentHandler struct {
	paymentUsecase  *usecases.PaymentUsecase
	historyUsecase  *usecases.HistoryUsecase
	customerUsecase *usecases.CustomerUsecase
}

func NewPaymentHandler(pu *usecases.PaymentUsecase, hu *usecases.HistoryUsecase, cu *usecases.CustomerUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase:  pu,
		historyUsecase:  hu,
		customerUsecase: cu,
	}
}

func (h *PaymentHandler) DoPayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromCustomerID string  `json:"from_customer_id"`
		ToCustomerID   string  `json:"to_customer_id"`
		Amount         float64 `json:"amount"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")
	log.Printf("token: %v", token)
	if token == "" {
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}

	customerID, err := h.customerUsecase.ValidateToken(token)
	log.Printf("customerID: %v", customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Printf("Before check")
	log.Printf("customerID: %v, fromCustomerID: %v", customerID, req.FromCustomerID)
	log.Print(customerID != req.FromCustomerID)
	if customerID != req.FromCustomerID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("After check")

	log.Printf("Before payment: %v", req.FromCustomerID)
	payment, err := h.paymentUsecase.DoPayment(req.FromCustomerID, req.ToCustomerID, req.Amount)
	log.Printf("After payment: %v", req.FromCustomerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.historyUsecase.LogAction("payment", req.FromCustomerID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(payment)
}
