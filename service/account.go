package service

import (
	"account-service-go-linkaja/model"

	"gorm.io/gorm"
)

type AccountService struct {
	DB *gorm.DB
}

func (as AccountService) Create(account *model.Accounts) error {
	// account.Balance += 10000

	err := as.DB.Create(account).Error
	return err
}
