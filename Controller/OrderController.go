package Controller

import (
	"github.com/gin-gonic/gin"
	"sunnybake/Model"
	"fmt"
	"net/http"
	time2 "time"
)

type OrderCon struct{

}

var OrderController OrderCon

func (orderCon OrderCon) View(ctx *gin.Context){
	db , err := Connect()
	defer db.Close()
	var transactionHeaders []Model.TransactionHeader
	if (err != nil){
		fmt.Println(err.Error())
	}else{
		db.Select(&transactionHeaders,"SELECT th.* FROM transactionheader th LEFT JOIN payment py ON th.id = py.transactionheader_id WHERE py.accountname not like '' and th.senddatetime =''")
		ctx.HTML(http.StatusOK,"order_view.html",gin.H{
			"transactionHeaders":transactionHeaders,
		})
	}
}
func (orderCon OrderCon) ViewDetail(ctx *gin.Context){
	transaction_header_id := ctx.Param("id")
	var transactionDetails [] Model.TransactionDetail
	db , err := Connect()
	defer db.Close()
	if(err != nil){
		fmt.Println(err.Error())
	}else{
		db.Select(&transactionDetails,"SELECT * FROM transactiondetail WHERE transaction_header_id = $1",transaction_header_id)
		var transactions []map[string]interface{}
		for _, transactionDetail:= range transactionDetails{
			var product Model.Product
			err:=db.Get(&product,"SELECT * FROM products WHERE id=$1",transactionDetail.Product_id)
			if(err != nil){
				fmt.Println(err.Error())
			}
			transaction:= map[string]interface{}{
				"transactionDetail":transactionDetail,
				"product":product,
			}
			transactions = append(transactions,transaction)
		}
		var transactionHeader Model.TransactionHeader
		err := db.Get(&transactionHeader,"SELECT * FROM transactionheader WHERE id = $1",transaction_header_id)
		if(err != nil){
			fmt.Println(err.Error())
		}

		ctx.HTML(http.StatusOK, "order_view_detail.html",gin.H{
			"transactions":transactions,
			"transactionHeader":transactionHeader,
		})
	}
}

func (orderCon OrderCon) ChangeStatus(ctx *gin.Context){
	currentTime := time2.Now()
	id:= ctx.Param("id")
	db, err := Connect()
	defer db.Close()
	if(err != nil){
		fmt.Println(err.Error())
	}else{
		db.MustExec("UPDATE transactionheader SET senddatetime = $1 WHERE id = $2",currentTime.Format(time2.UnixDate),id)
		ctx.Redirect(http.StatusSeeOther,"/order/view")
	}
}