package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sunnybake/Model"
	"fmt"
	"strconv"
	"database/sql"
	"net/smtp"
	"bytes"
	"html/template"
)

var TransactionController TransactionCon

type TransactionCon struct {
}

func (transaction TransactionCon) Index(ctx *gin.Context) {
	db, err := Connect()
	defer db.Close()
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
	defer db.Close()
	if err == nil {
		db.Select(&products, "SELECT * FROM products")
	}
	customer_name := ctx.PostForm("name")
	order_date := ctx.PostForm("date")
	address := ctx.PostForm("address")
	discount := ctx.PostForm("discount")
	price := ctx.PostForm("price")
	delivery_cost := ctx.PostForm("delivery")
	phoneNumber := ctx.PostForm("phoneNumber")

	var transactionHeader Model.TransactionHeader
	transactionHeader.CustomerName = customer_name
	transactionHeader.CustomerAddress = address
	transactionHeader.OrderDate = order_date
	transactionHeader.Discount = discount
	transactionHeader.Price = price
	transactionHeader.DeliveryCost = delivery_cost
	transactionHeader.SendDateTime = ""
	transactionHeader.PhoneNumber = phoneNumber
	_, err = db.NamedExec("INSERT INTO TransactionHeader(customername,customeraddress,discount,deliverycost,price,orderdate,senddatetime,phonenumber) VALUES (:customername, :customeraddress, :discount,:deliverycost,:price,:orderdate,:senddatetime,:phonenumber)", transactionHeader)
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
	ctx.Redirect(http.StatusSeeOther, "/transaction/view")

}
func (trasactionCon TransactionCon) View(ctx *gin.Context) {
	db, err := Connect()
	defer db.Close()
	if (err != nil) {
		fmt.Println(err.Error())
	} else {
		var transaction_headers []Model.TransactionHeader
		db.Select(&transaction_headers, "SELECT distinct th.id, th.customername, th.customeraddress, th.discount, th.deliverycost, th.price,th.orderdate,th.senddatetime FROM transactionheader th LEFT JOIN payment py ON th.id = py.transactionheader_id WHERE py.accountname like '' order By id Desc")
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
	defer db.Close()
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
	defer db.Close()
	if (err != nil) {
		fmt.Println(err.Error())
	} else {
		var transactionHeader Model.TransactionHeader;
		err = db.Get(&transactionHeader ,"SELECT * FROM transactionheader WHERE id = $1",transaction_header_id)
		if(err != nil){
			fmt.Println(err.Error())
		}
		qty :=0
		var transactionDetails []Model.TransactionDetail
		err = db.Select(&transactionDetails,"SELECT * FROM transactiondetail WHERE transaction_header_id = $1",transaction_header_id)
		if(err != nil){
			fmt.Println(err.Error())
		}else {
			var transactions []map[string]interface{}
			for _, transactionDetail := range transactionDetails {
				var product Model.Product
				err = db.Get(&product, "SELECT * FROM products WHERE id = $1", transactionDetail.Product_id)
				if (err != nil) {
					fmt.Println(err.Error())
				} else {
					subPrice ,_:=strconv.Atoi(product.Price)
					subPrice = subPrice * transactionDetail.Qty
					transaction := map[string]interface{}{
						"transactionDetail": transactionDetail,
						"product":           product,
						"subPrice":subPrice,
					}
					qty = qty + transactionDetail.Qty
					transactions = append(transactions, transaction)
				}
			}
			sendEmail(transactionHeader,transactions,qty)
			db.MustExec("UPDATE payment SET bankname=$1 , accountname=$2 , accountnumber=$3 , price=$4 , date=$5 WHERE transactionheader_id=$6", bank, name, number, price, date,transaction_header_id)
		}
	}
	ctx.Redirect(http.StatusSeeOther,"/order/view")
}
func (transactionCon TransactionCon) Delete(ctx *gin.Context){
	db , err := Connect()
	id := ctx.Param("id")
	defer  db.Close()
	if(err != nil){
		fmt.Println(err.Error())
	}else{
		db.MustExec("DELETE FROM transactionheader WHERE id=$1",id)
		db.MustExec("DELETE FROM transactiondetail WHERE transaction_header_id=$1",id)
	}
	ctx.Redirect(http.StatusSeeOther,"/transaction/view")
}

// sending Email
var auth smtp.Auth
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, "admsunnybake@gmail.com", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}



func sendEmail(th Model.TransactionHeader, ts []map[string]interface{} ,qty int){

	auth = smtp.PlainAuth("", "admsunnybake@gmail.com", "muarakarangpik", "smtp.gmail.com")
	/*templateData := struct{
		transactionHeader Model.TransactionHeader
		transactions []map[string]interface{}
	}{
		transactionHeader: th,
		transactions:  ts,
	}*/
	price ,_:= strconv.Atoi(th.Price)
	deliveryCost ,_:= strconv.Atoi(th.DeliveryCost)
	discount,_ := strconv.Atoi(th.Discount)
	totalPrice := price + deliveryCost- discount
	templateData:= gin.H{
		"transactionHeader":th,
		"transactions":ts,
		"totalPrice":totalPrice,
		"qty":qty,
	}

	r := NewRequest([]string{"Beatricedorothy@gmail.com"}, "Payment Successfully", "Sunny Bake Order")
	err := r.ParseTemplate("View/email.html", templateData)
	if  err == nil {
		ok, _ := r.SendEmail()
		fmt.Println(ok)
	}else{
		fmt.Println(err.Error())
	}




	// ini yang sebelumnya
	/*from := "admsunnybake@gmail.com"
	pass := "muarakarangpik"
	to := "Beatricedorothy@gmail.com"

	body := ""

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Payment Successfully \n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}*/

	//log.Print("sent, visit http://foobarbazz.mailinator.com")
}

