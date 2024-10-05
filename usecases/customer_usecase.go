package usecases

import (
	"errors"
	"log"
	"mnc-be-golang/models"
	"mnc-be-golang/repositories"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

type CustomerUsecase struct {
	customerRepo *repositories.CustomerRepository
}

func NewCustomerUsecase(repo *repositories.CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{customerRepo: repo}
}

type Claims struct {
	CustomerID string `json:"customer_id"`
	jwt.StandardClaims
}

func (u *CustomerUsecase) Login(email, password string) (*models.Customer, error) {
	customers, err := u.customerRepo.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, customer := range customers {
		if customer.Email == email && customer.Password == password {
			expirationTime := time.Now().Add(5 * time.Minute)
			claims := &Claims{
				CustomerID: customer.ID,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				return nil, err
			}

			customers[i].IsLoggedIn = true
			customers[i].Token = tokenString

			err = u.customerRepo.WriteAll(customers)
			if err != nil {
				return nil, err
			}

			return &customers[i], nil
		}
	}

	return nil, errors.New("invalid email or password")
}

func (u *CustomerUsecase) ValidateToken(tokenString string) (string, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &Claims{}
	log.Printf("tokenString after trim: %v", tokenString)
	log.Printf("claims: %v", claims)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	log.Printf("token: %v", token)
	log.Printf("err: %v", err)

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.CustomerID, nil
}

func (u *CustomerUsecase) Logout(customerID string) error {
	customers, err := u.customerRepo.ReadAll()
	if err != nil {
		return err
	}

	for i, customer := range customers {
		if customer.ID == customerID && customer.IsLoggedIn {
			customers[i].IsLoggedIn = false
			customers[i].Token = ""

			err = u.customerRepo.WriteAll(customers)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("customer not found or not logged in")
}
