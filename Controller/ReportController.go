package Controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/360EntSecGroup-Skylar/excelize"
	"net/http"
	"sunnybake/Model"
)

type ReportCon struct{

}

var ReportController ReportCon

func (reportCon ReportCon) View(ctx *gin.Context){
	ctx.HTML(http.StatusOK, "report.html",gin.H{

	})
}

func (reportCon ReportCon) Generate(ctx *gin.Context){
	var transactionHeaders []Model.TransactionHeader
	start_date := ctx.PostForm("start_date")
	end_date := ctx.PostForm("end_date")
	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "Order Date")
	xlsx.SetCellValue("Sheet1", "B1", "Name")
	xlsx.SetCellValue("Sheet1", "C1", "Price")
	xlsx.SetCellValue("Sheet1", "D1", "Delivery Cost")
	xlsx.SetCellValue("Sheet1", "E1", "Discount")
	xlsx.SetCellValue("Sheet1", "F1", "Product")
	xlsx.SetCellValue("Sheet1", "G1", "Price Per Item")
	xlsx.SetCellValue("Sheet1", "H1", "Quantity")
	db, err := Connect()
	// err = db.Get(&transactionHeader, "SELECT id,customername,customeraddress,discount,deliverycost,price,orderdate,phonenumber, delivery,concat(substring(senddatetime, 1,position (':' in senddatetime )-3),' ',right(senddatetime ,4),' ',substring(senddatetime , position (':' in senddatetime )-2 , 8)) as senddatetime FROM transactionheader WHERE id = $1", id)
	err = db.Select(&transactionHeaders, "SELECT id,customername,customeraddress,discount,deliverycost,price,orderdate,phonenumber, delivery FROM transactionheader WHERE DATE(orderdate) BETWEEN DATE($1) AND DATE($2)", start_date, end_date)
	if err != nil {
		fmt.Println(err.Error())
	}else{
		var row int =2 
		for _, transactionHeader := range transactionHeaders{
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d",row), transactionHeader.OrderDate)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d",row), transactionHeader.CustomerName)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d",row), transactionHeader.Price)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d",row), transactionHeader.DeliveryCost)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d",row), transactionHeader.Discount)
			var transactionDetails []Model.TransactionDetail
			err = db.Select(&transactionDetails, "SELECT * FROM transactiondetail WHERE transaction_header_id = $1", transactionHeader.Id)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				for _, transactionDetail := range transactionDetails {
					var product Model.Product
					err = db.Get(&product, "SELECT * FROM products WHERE id = $1", transactionDetail.Product_id)
					if err != nil {
						fmt.Println(err.Error())
					} else {
						xlsx.SetCellValue("Sheet1", fmt.Sprintf("F%d",row), product.Name)
						xlsx.SetCellValue("Sheet1", fmt.Sprintf("G%d",row), product.Price)
						xlsx.SetCellValue("Sheet1", fmt.Sprintf("H%d",row), transactionDetail.Qty)
						row++
					}
				}
			}
		}
	}
	
	//_ = xlsx.SaveAs("./aaa.xlsx")
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+"report.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")
	_ = xlsx.Write(ctx.Writer)

}
