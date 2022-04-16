package model

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

//receiver function
func (Customers) TableName() string {
	return "customers"
}

//receiver function
func (Accounts) TableName() string {
	return "accounts"
}
