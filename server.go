package main

import (
	"account-service-go-linkaja/db"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

//struct for account
type Accounts struct {
	AccountNumber int    `json:"account_number"`
	CustomerName  string `json:"customer_number"`
	Balance       int    `json:"balance"`
}

//receiver function
func (Accounts) TableName() string {
	return "users"
}

func main() {
	route := echo.New()
	db := db.Init()
	route.POST("account/create_account", func(c echo.Context) error {
		account := new(Accounts)
		c.Bind(account)

		err := db.Create(account).Error
		if err != nil {
			fmt.Println("Error created")
			return c.JSON(http.StatusBadRequest, err)
		}
		response := struct {
			Message string
			Data    Accounts
		}{
			Message: "New account has been created successfully",
			Data:    *account,
		}

		return c.JSON(http.StatusOK, response)
	})
}
