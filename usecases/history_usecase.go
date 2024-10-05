package usecases

import (
	"mnc-be-golang/models"
	"mnc-be-golang/repositories"
	"time"
)

type HistoryUsecase struct {
	repo *repositories.HistoryRepository
}

func NewHistoryUsecase(repo *repositories.HistoryRepository) *HistoryUsecase {
	return &HistoryUsecase{repo: repo}
}

func (u *HistoryUsecase) LogAction(action, customerID string) error {
	history := models.History{
		ID:         time.Now().Format(time.RFC3339),
		Email:      customerID,
		Password:   customerID,
		IsLoggedIn: true,
	}

	return u.repo.WriteHistory(history)
}
