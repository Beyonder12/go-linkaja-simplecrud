package test

import (
	"account-service-go-linkaja/db"
	"account-service-go-linkaja/model"
	"account-service-go-linkaja/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

// We need this function to kick off the test suite, otherwise
// "go test" won't know about our tests
func TestCreateAccount(t *testing.T) {
	mockDb := db.Init()
	accountService := service.AccountService{
		DB: mockDb,
	}

	data := model.Accounts{
		AccountNumber:  1000,
		CustomerNumber: 1,
		Balance:        20000,
	}

	err := accountService.Create(&data)
	assert.Nil(t, err)
	assert.Equal(t, 20000, data.Balance)

	// result := model.Accounts{
	// 	ID: data.ID,
	// }
	// err = db.First(&result).Error
	// assert.Nil(t, err)

	// assert.Equal(t, 1000, result.AccountNumber)
	//go test -v ./test
}
