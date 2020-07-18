package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sunnybake/Model"
	"fmt"
	"strings"
)

var ProductController ProductControl

type ProductControl struct {
}

func (productController ProductControl) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "insert_product.html", gin.H{

	})

}

func (productController ProductControl) Insert(ctx *gin.Context) {
	var product Model.Product;
	product.Name = ctx.PostForm("name")
	product.Name = strings.Trim(product.Name, " ")
	product.Price = ctx.PostForm("price")
	product.Description = ctx.PostForm("description")
	file, _ := ctx.FormFile("photo")

	//fmt.Println(file.Filename)
	product.Url = strings.Replace(product.Name, " ", "_", -1) + "_" + strings.Replace(file.Filename, " ", "_", -1)
	//fmt.Println(product.Url)
	ctx.SaveUploadedFile(file, "static/images/"+product.Url)

	cloudinary := Model.Create("148941686835669","hj-ZYCdO6jUpiwunoh2Hu9yUgO4","sunnybake")
	options := map[string]string{
		"public_id": product.Url,
	}
	images,er:=cloudinary.Upload("https://bttr.herokuapp.com/asset/images/"+product.Url,options)
	if(er == nil){
		product.Url = images.Url
	}else{
		fmt.Println(er)
	}



	db, err := Connect()
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	db.NamedExec("INSERT INTO products(name,price,url,description) VALUES(:name,:price,:url,:description)", product)
	ctx.Redirect(http.StatusSeeOther, "/")
}

func (productController ProductControl) ViewUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	//fmt.Println(id)
	db, err := Connect()
	defer db.Close()
	if (err != nil) {
		fmt.Println(err.Error())
	} else {
		var product Model.Product
		db.Get(&product, "SELECT * FROM products WHERE id=$1", id)
		//fmt.Println(product)
		ctx.HTML(http.StatusOK, "update_product.html", gin.H{
			"product": product,
		})

	}
	db.Close()

}

func (productControl ProductControl) Update(ctx *gin.Context) {
	var product Model.Product;
	db, err := Connect()
	defer db.Close()
	id:= ctx.Param("id")
	if (err != nil) {
		fmt.Println(err.Error())
	} else {
		db.Get(&product, "SELECT * FROM products WHERE id=$1", id)
		product.Name = ctx.PostForm("name")
		product.Name = strings.Trim(product.Name, " ")
		product.Price = ctx.PostForm("price")
		product.Description = ctx.PostForm("description")
		file, _ := ctx.FormFile("photo")
		//fmt.Println(file.Filename)
		if (file.Filename != "") {
			product.Url = strings.Replace(product.Name, " ", "_", -1) + "_" + strings.Replace(file.Filename, " ", "_", -1)
			//fmt.Println(product.Url)
			ctx.SaveUploadedFile(file, "static/images/"+product.Url)

			cloudinary := Model.Create("148941686835669","hj-ZYCdO6jUpiwunoh2Hu9yUgO4","sunnybake")
			options := map[string]string{
				"public_id": product.Url,
			}
			images,er:=cloudinary.Upload("https://bttr.herokuapp.com/asset/images/"+product.Url,options)
			if(er == nil){
				product.Url = images.Url
			}else{
				fmt.Println(er)
			}
		}

		db.NamedExec("UPDATE products SET name=:name , price=:price , description=:description,url=:url WHERE id=:id",product)
		ctx.Redirect(http.StatusSeeOther,"/")
	}
}

func (productControl ProductControl) Delete(ctx * gin.Context){
	id := ctx.Param("id")
	db , err := Connect()
	defer db.Close()
	if(err != nil){
		fmt.Println(err.Error())
	}else{
		db.MustExec("DELETE FROM products WHERE id=$1",id)
		ctx.Redirect(http.StatusSeeOther,"/")
	}
}