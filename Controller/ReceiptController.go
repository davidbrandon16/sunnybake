package Controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sunnybake/Model"
)

type ReceiptCon struct {
}

var ReceiptController ReceiptCon

func (receiptCon ReceiptCon) View(ctx *gin.Context) {
	db, err := Connect()
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var transactionHeaders []Model.TransactionHeader
		err := db.Select(&transactionHeaders, "SELECT * FROM transactionheader ORDER BY orderdate::date DESC LIMIT 10")
		if err != nil {
			fmt.Println(err.Error())
		} else {
			ctx.HTML(http.StatusOK, "view_receipt.html", gin.H{
				"transactionHeaders": transactionHeaders,
			})
		}
	}
}

func (receiptCon ReceiptCon) Print(ctx *gin.Context) {
	id := ctx.Param("id")
	db, err := Connect()
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var transactionHeader Model.TransactionHeader
		// err = db.Get(&transactionHeader, "SELECT id,customername,customeraddress,discount,deliverycost,price,orderdate,phonenumber, delivery,concat(substring(senddatetime, 1,position (':' in senddatetime )-3),' ',right(senddatetime ,4),' ',substring(senddatetime , position (':' in senddatetime )-2 , 8)) as senddatetime FROM transactionheader WHERE id = $1", id)
		err = db.Get(&transactionHeader, "SELECT id,customername,customeraddress,discount,deliverycost,price,orderdate,phonenumber, delivery FROM transactionheader WHERE id = $1", id)
		if err != nil {
			fmt.Println(err.Error())
		}
		var transactionDetails []Model.TransactionDetail
		err = db.Select(&transactionDetails, "SELECT * FROM transactiondetail WHERE transaction_header_id = $1", id)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			var transactions []map[string]interface{}
			for _, transactionDetail := range transactionDetails {
				var product Model.Product
				err = db.Get(&product, "SELECT * FROM products WHERE id = $1", transactionDetail.Product_id)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					transaction := map[string]interface{}{
						"transactionDetail": transactionDetail,
						"product":           product,
					}
					transactions = append(transactions, transaction)
				}
			}
			ctx.HTML(http.StatusOK, "print_receipt.html", gin.H{
				"transactionHeader": transactionHeader,
				"transactions":      transactions,
			})
		}
	}
}

func (reciptCon ReceiptCon) ViewPage(ctx *gin.Context) {
	var total_data int
	db, err := Connect()
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		//var transaction_headers[] Model.TransactionHeader
		db.Get(&total_data, "SELECT count(id) FROM transactionheader ")
		var page int
		page, err = strconv.Atoi(ctx.Param("page"))
		data := page * 10
		if data < total_data {
			var until int
			if total_data-data < 10 {
				until = total_data - data
			} else {
				until = 10
			}
			var transactionHeaders []Model.TransactionHeader
			db.Select(&transactionHeaders, "SELECT * FROM transactionheader ORDER BY orderdate Desc LIMIT $1 OFFSET $2", until, data+1)
			ctx.JSON(http.StatusOK, gin.H{
				"transactionHeaders": transactionHeaders,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"transactionHeaders": nil,
		})

	}

}
