package usecases

import (
	"errors"
	"mnc-be-golang/models"
	"mnc-be-golang/repositories"
	"time"
)

type PaymentUsecase struct {
	customerRepo *repositories.CustomerRepository
	paymentRepo  *repositories.PaymentRepository
}

func NewPaymentUsecase(customerRepo *repositories.CustomerRepository, paymentRepo *repositories.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{
		customerRepo: customerRepo,
		paymentRepo:  paymentRepo,
	}
}

func (u *PaymentUsecase) DoPayment(fromCustomerID, toCustomerID string, amount float64) (*models.Payment, error) {
	customers, err := u.customerRepo.ReadAll()
	if err != nil {
		return nil, err
	}

	var fromCustomer *models.Customer
	var toCustomer *models.Customer

	for i, customer := range customers {
		if customer.ID == fromCustomerID && customer.IsLoggedIn {
			fromCustomer = &customers[i]
		}
		if customer.ID == toCustomerID {
			toCustomer = &customers[i]
		}
	}

	if fromCustomer == nil {
		return nil, errors.New("customer is not logged in")
	}

	if toCustomer == nil {
		return nil, errors.New("recipient customer not found")
	}

	if fromCustomer.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	fromCustomer.Balance -= amount
	toCustomer.Balance += amount

	err = u.customerRepo.UpdateCustomerBalance(fromCustomer.ID, fromCustomer.Balance)
	if err != nil {
		return nil, err
	}

	err = u.customerRepo.UpdateCustomerBalance(toCustomer.ID, toCustomer.Balance)
	if err != nil {
		return nil, err
	}

	payment := models.Payment{
		FromCustomerID: fromCustomerID,
		ToCustomerID:   toCustomerID,
		Amount:         amount,
		Timestamp:      time.Now().Format(time.RFC3339),
	}

	err = u.paymentRepo.WritePayment(payment)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}
