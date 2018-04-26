package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../Model"
	"fmt"
	"strconv"
	"database/sql"
)

var TransactionController TransactionCon

type TransactionCon struct {
}

func (transaction TransactionCon) Index(ctx *gin.Context) {
	db, err := Connect()
	var products []Model.Product
	if (err == nil) {
		db.Select(&products, "SELECT * FROM products")
	}

	ctx.HTML(http.StatusOK, "insert_transaction.html", gin.H{
		"products": products,
	})

}
func (transaction TransactionCon) Insert(ctx *gin.Context) {
	var products []Model.Product
	db, err := Connect()
	if err == nil {
		db.Select(&products, "SELECT * FROM products")
	}
	customer_name := ctx.PostForm("name")
	order_date := ctx.PostForm("date")
	address := ctx.PostForm("address")
	discount := ctx.PostForm("discount")
	price := ctx.PostForm("price")
	delivery_cost := ctx.PostForm("delivery")

	var transactionHeader Model.TransactionHeader
	transactionHeader.CustomerName = customer_name
	transactionHeader.CustomerAddress = address
	transactionHeader.OrderDate = order_date
	transactionHeader.Discount = discount
	transactionHeader.Price = price
	transactionHeader.DeliveryCost = delivery_cost
	transactionHeader.SendDateTime = ""
	_, err = db.NamedExec("INSERT INTO TransactionHeader(customername,customeraddress,discount,deliverycost,price,orderdate,senddatetime) VALUES (:customername, :customeraddress, :discount,:deliverycost,:price,:orderdate,:senddatetime)", transactionHeader)
	if (err != nil) {
		fmt.Println(err.Error())
	}
	err = db.Get(&transactionHeader, "SELECT * FROM transactionheader ORDER BY Id Desc limit 1")
	if (err != nil) {
		fmt.Println(err.Error())
	}
	for _, product := range products {
		qty := ctx.PostForm("qty" + strconv.Itoa(product.Id))
		if (qty != "") {
			var transactionDetail Model.TransactionDetail
			transactionDetail.Product_id = product.Id
			transactionDetail.Qty, _ = strconv.Atoi(qty)
			transactionDetail.Transaction_Header_id = transactionHeader.Id
			db.NamedExec("INSERT INTO transactionDetail(product_id,qty,transaction_header_id) VALUES (:product_id,:qty,:transaction_header_id)", transactionDetail)
		}
	}
	var payment Model.Payment
	payment.TransactionHeader_id = transactionHeader.Id
	payment.AccountName = ""
	payment.BankName = ""
	payment.AccountNumber = ""
	payment.Date = ""
	db.NamedExec("INSERT INTO payment(bankname,accountname,accountnumber,date,transactionheader_id) VALUES (:bankname,:accountname,:accountnumber,:date,:transactionheader_id)", payment)
	ctx.Redirect(http.StatusSeeOther, "/")

}
func (trasactionCon TransactionCon) View(ctx *gin.Context) {
	db, err := Connect()
	if (err != nil) {
		fmt.Println(err.Error())
	} else {
		var transaction_headers []Model.TransactionHeader
		db.Select(&transaction_headers, "SELECT th.id, th.customername, th.customeraddress, th.discount, th.deliverycost, th.price,th.orderdate,th.senddatetime FROM transactionheader th LEFT JOIN payment py ON th.id = py.transactionheader_id WHERE py.accountname like ''")
		var obj_transaction_headers []interface{}
		for _, transaction_header := range transaction_headers {
			var transaction_details []Model.TransactionDetail
			err := db.Select(&transaction_details, "SELECT * FROM transactiondetail where transaction_header_id = $1", transaction_header.Id)
			if (err != nil) {
				fmt.Println(err.Error())
			}
			/*var obj_transaction_details []interface{}
			for _, transaction_detail := range  transaction_details {

				var product Model.Product
				err := db.Get(&product,"SELECT * FROM products WHERE id = $1",transaction_detail.Product_id)
				if(err != nil ){
					fmt.Println(err.Error())
				}
				obj_transaction_detail := map[string] interface{}{
					"product" : product,
				}
				fmt.Println(obj_transaction_detail)
				obj_transaction_details = append(obj_transaction_details,obj_transaction_detail)
			}*/
			obj_transaction_header := map[string]interface{}{
				"transactionHeader": transaction_header,
			}
			obj_transaction_headers = append(obj_transaction_headers, obj_transaction_header)
		}
		ctx.HTML(http.StatusOK, "view_transaction.html", gin.H{
			"transactions": obj_transaction_headers,
		})
	}
}

func (transactionCon TransactionCon) ChangePayment(ctx *gin.Context) {
	id := ctx.Param("id")

	db, err := Connect()
	if (err == nil) {
		var transactionHeader Model.TransactionHeader
		err := db.Get(&transactionHeader, "SELECT th.* FROM transactionheader th JOIN payment py ON py.transactionheader_id = th.id WHERE th.id = $1 and py.accountname like ''", id)
		if (err != nil) {
			if(err == sql.ErrNoRows){
				ctx.Redirect(http.StatusSeeOther,"/")
			}
			fmt.Println(err.Error())
		} else {
			ctx.HTML(http.StatusOK, "change_payment.html", gin.H{
				"transactionHeader": transactionHeader,
			})
		}
	}
}

func (transactionCon TransactionCon) InsertPayment(ctx *gin.Context) {
	transaction_header_id := ctx.Param("id")
	bank := ctx.PostForm("bank")
	name := ctx.PostForm("accountName")
	number := ctx.PostForm("accountNumber")
	price := ctx.PostForm("price")
	date := ctx.PostForm("date")
	db, err := Connect()
	if (err != nil) {
		fmt.Println(err.Error())
	} else {
		db.MustExec("UPDATE payment SET bankname=$1 , accountname=$2 , accountnumber=$3 , price=$4 , date=$5 WHERE transactionheader_id=$6", bank, name, number, price, date,transaction_header_id)
	}
	ctx.Redirect(http.StatusSeeOther,"/transaction/view")
}
