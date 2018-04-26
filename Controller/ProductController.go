package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sunnybake/Model"
	"fmt"
	"strings"
)

var ProductController ProductControl

type ProductControl struct{
}

func (productController ProductControl) Index(ctx *gin.Context){
	ctx.HTML(http.StatusOK,"insert_product.html",gin.H{

	})


}

func (productController ProductControl) Insert(ctx *gin.Context){
	var product Model.Product;
	product.Name= ctx.PostForm("name")
	product.Name= strings.Trim(product.Name," ")
	product.Price = ctx.PostForm("price")
	product.Description = ctx.PostForm("description")
	file ,_:= ctx.FormFile("photo")

	fmt.Println(file.Filename)
	product.Url = strings.Replace(product.Name," ","_",-1)+"_"+strings.Replace(file.Filename," ","_",-1)
	fmt.Println(product.Url)
	ctx.SaveUploadedFile(file,"static/images/"+product.Url)


	db ,err := Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	db.NamedExec("INSERT INTO products(name,price,url,description) VALUES(:name,:price,:url,:description)",product)
	db.Close()
	ctx.Redirect(http.StatusSeeOther,"/")
}