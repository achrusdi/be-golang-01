package main

import (
	"mnc-be-golang/handlers"
	"mnc-be-golang/repositories"
	"mnc-be-golang/usecases"
	"net/http"
)

func main() {
	customerRepo := repositories.NewCustomerRepository("storages/customer.json")
	historyRepo := repositories.NewHistoryRepository("storages/history.json")
	paymentRepo := repositories.NewPaymentRepository("storages/payment.json")

	customerUsecase := usecases.NewCustomerUsecase(customerRepo)
	historyUsecase := usecases.NewHistoryUsecase(historyRepo)
	paymentUsecase := usecases.NewPaymentUsecase(customerRepo, paymentRepo)

	customerHandler := handlers.NewCustomerHandler(customerUsecase, historyUsecase)
	paymentHandler := handlers.NewPaymentHandler(paymentUsecase, historyUsecase, customerUsecase)

	http.HandleFunc("/login", customerHandler.Login)
	http.HandleFunc("/logout", customerHandler.Logout)
	http.HandleFunc("/payment", paymentHandler.DoPayment)

	http.ListenAndServe(":8080", nil)
}
