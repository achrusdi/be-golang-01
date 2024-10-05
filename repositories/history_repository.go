package repositories

import (
	"encoding/json"
	"mnc-be-golang/models"
	"os"
)

type HistoryRepository struct {
	filePath string
}

func NewHistoryRepository(filePath string) *HistoryRepository {
	return &HistoryRepository{filePath: filePath}
}

func (r *HistoryRepository) WriteHistory(history models.History) error {
	histories, err := r.ReadAll()
	if err != nil {
		return err
	}

	histories = append(histories, history)

	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(histories)
}

func (r *HistoryRepository) ReadAll() ([]models.History, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var histories []models.History
	err = json.NewDecoder(file).Decode(&histories)
	if err != nil {
		return nil, err
	}
	return histories, nil
}
