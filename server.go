package main

import (
	"account-service-go-linkaja/db"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

//struct for customer
type Customers struct {
	ID       int      `json:"customer_number"`
	Name     string   `json:"name"`
	Accounts Accounts `gorm:"foreignKey:CustomerNumber"`
}

//struct for account
type Accounts struct {
	ID             int `json:"id"`
	AccountNumber  int `json:"account_number"`
	CustomerNumber int `json:"customer_number"`
	Balance        int `json:"balance"`
}

//struct dto for trasfer
type TransferDto struct {
	ToAccountNumber int `json:"to_account_number"`
	Amount          int `json:"amount"`
}

type AccountResponseDto struct {
	AccountNumber int    `json:"account_number"`
	CustomerName  string `json:"customer_name"`
	Balance       int    `json:"balance"`
}

//receiver function
func (Customers) TableName() string {
	return "customers"
}

//receiver function
func (Accounts) TableName() string {
	return "accounts"
}

func main() {
	route := echo.New()
	db := db.Init()
	route.POST("account/", func(c echo.Context) error {
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

	route.GET("account/", func(c echo.Context) error {

		var customers []Customers
		db.Joins("Accounts").Find(&customers)

		var results []AccountResponseDto
		for _, cus := range customers {
			a := AccountResponseDto{
				AccountNumber: cus.Accounts.AccountNumber,
				CustomerName:  cus.Name,
				Balance:       cus.Accounts.Balance,
			}

			results = append(results, a)
		}

		response := struct {
			Message string
			Data    []AccountResponseDto
		}{
			Message: "Successfully fetch all account's data",
			Data:    results,
		}

		return c.JSON(http.StatusOK, response)
	})

	route.GET("account/:account_number", func(c echo.Context) error {
		account := new(Accounts)
		account.AccountNumber, _ = strconv.Atoi(c.Param("account_number"))

		err := db.Where("account_number", account.AccountNumber).First(&account).Error
		if err != nil {
			fmt.Println("account number not found")
			return c.JSON(http.StatusBadRequest, err)
		}
		response := struct {
			Message string
			Data    Accounts
		}{
			Message: "Succesfully found an account",
			Data:    *account,
		}
		return c.JSON(http.StatusOK, response)
	})

	route.POST("account/:from_account_number/transfer", func(c echo.Context) error {
		//get fromaccount from param respectively
		fromAccount := new(Accounts)
		fromAccount.AccountNumber, _ = strconv.Atoi(c.Param("from_account_number"))

		errFrom := db.Where("account_number = ?", fromAccount.AccountNumber).First(&fromAccount).Error
		if errFrom != nil {
			fmt.Println("from account number is not found")
			return c.JSON(http.StatusNotFound, errFrom)
		}

		// get account using body request toAccountNumber
		transferDto := new(TransferDto)
		c.Bind(transferDto)

		var toAccount Accounts

		errTo := db.Where("account_number = ?", transferDto.ToAccountNumber).Find(&toAccount).Error
		if errFrom != nil {
			fmt.Println("error find account")
			return c.JSON(http.StatusNotFound, errTo)
		}

		fromAccount.Balance = fromAccount.Balance - transferDto.Amount

		errFrom = db.Save(fromAccount).Error
		if errFrom != nil {
			fmt.Println("error save from account")
			return c.JSON(http.StatusBadRequest, errFrom)
		}

		toAccount.Balance = toAccount.Balance + transferDto.Amount

		errTo = db.Save(toAccount).Error
		if errTo != nil {
			fmt.Println("error save to account")
			return c.JSON(http.StatusBadRequest, errTo)
		}

		// update balance
		response := struct {
			Message     string
			FromAccount Accounts
			ToAccount   Accounts
		}{
			Message:     "Transaction success",
			FromAccount: *fromAccount,
			ToAccount:   toAccount,
		}

		return c.JSON(http.StatusAccepted, response)

	})

	route.Start(":9999")

}
