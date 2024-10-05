package handlers

import (
	"encoding/json"
	"mnc-be-golang/usecases"
	"net/http"
	"strings"
)

type CustomerHandler struct {
	customerUsecase *usecases.CustomerUsecase
	historyUsecase  *usecases.HistoryUsecase
}

func NewCustomerHandler(cu *usecases.CustomerUsecase, hu *usecases.HistoryUsecase) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: cu,
		historyUsecase:  hu,
	}
}

func (h *CustomerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	customer, err := h.customerUsecase.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.historyUsecase.LogAction("login", customer.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	customerID, err := h.customerUsecase.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	err = h.customerUsecase.Logout(customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.historyUsecase.LogAction("logout", customerID)

	response := map[string]string{
		"message":     "Logout successful",
		"customer_id": customerID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
