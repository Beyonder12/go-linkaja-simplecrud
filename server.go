package main

import (
	"account-service-go-linkaja/db"
	"account-service-go-linkaja/model"
	"account-service-go-linkaja/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

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

func main() {
	route := echo.New()
	db := db.Init()

	customerService := service.CustomerService{
		Db: db,
	}
	
	accountService := service.AccountService{
		DB: db,
	}
	

	route.POST("customer/", func(c echo.Context) error {
		customer := new(model.Customers)
		c.Bind(customer)

		err := customerService.Create(customer)
		if err != nil {
			fmt.Println("Error created")
			return c.JSON(http.StatusBadRequest, err)
		}
		response := struct {
			Message string
			Data    model.Customers
		}{
			Message: "New customer has been created successfully",
			Data: *customer
		}
	})

	route.POST("account/", func(c echo.Context) error {
		account := new(model.Accounts)
		c.Bind(account)

		err := accountService.Create(account)
		if err != nil {
			fmt.Println("Error created")
			return c.JSON(http.StatusBadRequest, err)
		}
		response := struct {
			Message string
			Data    model.Accounts
		}{
			Message: "New account has been created successfully",
			Data:    *account,
		}

		return c.JSON(http.StatusOK, response)
	})

	route.GET("account/", func(c echo.Context) error {

		var customers []model.Customers
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
		account := new(model.Accounts)
		account.AccountNumber, _ = strconv.Atoi(c.Param("account_number"))

		err := db.Where("account_number", account.AccountNumber).First(&account).Error
		if err != nil {
			fmt.Println("account number not found")
			return c.JSON(http.StatusBadRequest, err)
		}
		response := struct {
			Message string
			Data    model.Accounts
		}{
			Message: "Succesfully found an account",
			Data:    *account,
		}
		return c.JSON(http.StatusOK, response)
	})

	route.POST("account/:from_account_number/transfer", func(c echo.Context) error {
		//get fromaccount from param respectively
		fromAccount := new(model.Accounts)
		fromAccount.AccountNumber, _ = strconv.Atoi(c.Param("from_account_number"))

		errFrom := db.Where("account_number = ?", fromAccount.AccountNumber).First(&fromAccount).Error
		if errFrom != nil {
			fmt.Println("from account number is not found")
			return c.JSON(http.StatusNotFound, errFrom)
		}

		// get account using body request toAccountNumber
		transferDto := new(TransferDto)
		c.Bind(transferDto)

		var toAccount model.Accounts

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
			FromAccount model.Accounts
			ToAccount   model.Accounts
		}{
			Message:     "Transaction success",
			FromAccount: *fromAccount,
			ToAccount:   toAccount,
		}

		return c.JSON(http.StatusAccepted, response)

	})

	route.Start(":9999")

}
