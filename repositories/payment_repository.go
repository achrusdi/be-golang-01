package repositories

import (
	"encoding/json"
	"mnc-be-golang/models"
	"os"
)

type PaymentRepository struct {
	filePath string
}

func NewPaymentRepository(filePath string) *PaymentRepository {
	return &PaymentRepository{filePath: filePath}
}

func (r *PaymentRepository) ReadAll() ([]models.Payment, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var payments []models.Payment
	err = json.NewDecoder(file).Decode(&payments)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) WriteAll(payments []models.Payment) error {
	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(payments)
}

func (r *PaymentRepository) WritePayment(payment models.Payment) error {
	payments, err := r.ReadAll()
	if err != nil {
		return err
	}

	payments = append(payments, payment)

	return r.WriteAll(payments)
}
