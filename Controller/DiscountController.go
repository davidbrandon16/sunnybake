package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sunnybake/Model"
	"fmt"
	"strconv"
)

var DiscountController DiscountControl

type DiscountControl struct {
}

func (discountControl DiscountControl) Index(ctx *gin.Context){
	ctx.HTML(http.StatusOK,"insert_discount.html",gin.H{

	})
}


func (discountControl DiscountControl) Insert(ctx * gin.Context){
	name:= ctx.PostForm("name")
	Type:= ctx.PostForm("type")
	price := ctx.PostForm("price")
	var discount Model.Discount
	discount.Name= name
	discount.Type = Type
	discount.Price = price
	db , err := Connect()
	defer db.Close()
	if(err != nil){
		fmt.Println(err.Error())
	}else{
		db.NamedExec("INSERT INTO discount(name,type,price) VALUES (:name,:type,:price)",discount)
	}
	ctx.Redirect(http.StatusSeeOther,"/discount/view")
}

func (discountControl DiscountControl) Update(ctx *gin.Context){
	id := ctx.Param("id")
	var discount Model.Discount
	db , err := Connect()
	defer  db.Close()
	if(err != nil){
		fmt.Println(err.Error())
	}else{
		db.Get(&discount , "SELECT * FROM discount WHERE id = $1",id)
		ctx.HTML(http.StatusOK,"update_discount.html",gin.H{
			"discount":discount,
		})
	}
}

func (discountControl DiscountControl) UpdateData(ctx *gin.Context){
	id := ctx.Param("id")
	name:= ctx.PostForm("name")
	Type:= ctx.PostForm("type")
	price := ctx.PostForm("price")
	var discount Model.Discount
	discount.Name= name
	discount.Type = Type
	discount.Price = price
	discount.Id ,_= strconv.Atoi(id)
	db , err := Connect()
	defer db.Close()
	if(err != nil){
		fmt.Println(err.Error())
	}else{
		db.NamedExec("UPDATE discount SET name=:name,type = :type,price =:price WHERE id = :id",discount)
	}
	ctx.Redirect(http.StatusSeeOther,"/discount/view")
}

func (DiscountControl DiscountControl) View(ctx *gin.Context){
	db, err := Connect()
	if (err != nil){
		fmt.Println(err.Error())
	}else{
		var discounts [] Model.Discount
		err := db.Select(&discounts,"SELECT * FROM discount")
		if(err != nil) {
			discounts = nil
		}
		ctx.HTML(http.StatusOK,"view_discount.html",gin.H{
			"discounts":discounts,
		})
	}
}
func (discountControl DiscountControl) Delete(ctx *gin.Context){
	id := ctx.Param("id")
	db , err := Connect()
	defer  db.Close()
	if(err != nil){
		fmt.Println(err.Error())
	}else{
		db.MustExec("DELETE FROM discount WHERE id=$1",id)
		ctx.Redirect(http.StatusSeeOther,"/discount/view")

	}
}