package repositories

import (
	"encoding/json"
	"mnc-be-golang/models"
	"os"
)

type CustomerRepository struct {
	filePath string
}

func NewCustomerRepository(filePath string) *CustomerRepository {
	return &CustomerRepository{filePath: filePath}
}

func (r *CustomerRepository) ReadAll() ([]models.Customer, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var customers []models.Customer
	err = json.NewDecoder(file).Decode(&customers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomerRepository) WriteAll(customers []models.Customer) error {
	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(customers)
}

func (r *CustomerRepository) UpdateCustomerBalance(customerID string, newBalance float64) error {
	customers, err := r.ReadAll()
	if err != nil {
		return err
	}

	for i, customer := range customers {
		if customer.ID == customerID {
			customers[i].Balance = newBalance
		}
	}

	return r.WriteAll(customers)
}
