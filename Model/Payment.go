package Model


type Payment struct {
	Id int
	TransactionHeader_id int
	BankName string
	AccountNumber string
	AccountName string
	Date string
	Price float64
}