package Controller

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"sunnybake/Model"
	"net/http"
)

type ReceiptCon struct {

}

var ReceiptController ReceiptCon

func (receiptCon ReceiptCon) View(ctx * gin.Context){
	db , err := Connect()
	defer db.Close()
	if(err != nil){
		fmt.Println(err.Error())
	}else{
		var transactionHeaders []Model.TransactionHeader
		err:=db.Select(&transactionHeaders,"SELECT * FROM transactionheader ")
		if(err != nil){
			fmt.Println(err.Error())
		}else{
			ctx.HTML(http.StatusOK,"view_receipt.html",gin.H{
				"transactionHeaders":transactionHeaders,
			})
		}
	}
}

func ( receiptCon ReceiptCon) Print(ctx * gin.Context){
	id:= ctx.Param("id")
	db, err:= Connect()
	defer db.Close()
	if(err != nil){
		fmt.Println(err.Error())
	}else{
		var transactionHeader Model.TransactionHeader
		err = db.Get(&transactionHeader ,"SELECT * FROM transactionheader WHERE id = $1",id)
		if(err != nil){
			fmt.Println(err.Error())
		}
		var transactionDetails []Model.TransactionDetail
		err = db.Select(&transactionDetails,"SELECT * FROM transactiondetail WHERE transaction_header_id = $1",id)
		if(err != nil){
			fmt.Println(err.Error())
		}else{
			var transactions []map[string] interface{}
			for _,transactionDetail := range transactionDetails{
				var product Model.Product
				err = db.Get(&product,"SELECT * FROM products WHERE id = $1",transactionDetail.Product_id)
				if(err != nil){
					fmt.Println(err.Error())
				}else {
					transaction := map[string]interface{}{
						"transactionDetail": transactionDetail,
						"product":product,
					}
					transactions = append(transactions,transaction)
				}
			}
			ctx.HTML(http.StatusOK,"print_receipt.html",gin.H{
				"transactionHeader":transactionHeader,
				"transactions":transactions,
			})
		}
	}
}