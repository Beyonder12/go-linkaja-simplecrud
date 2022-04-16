package service

import (
	"account-service-go-linkaja/model"

	"gorm.io/gorm"
)

type CustomerService struct {
	DB *gorm.DB
}

func (cs CustomerService) Create(customer *model.Customers) error {
	err := cs.DB.Create(customer).Error

	return err
}
