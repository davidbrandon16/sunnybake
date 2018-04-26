package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"sunnybake/Model"
)
var HomeController Home

type Home struct{

}

func (homeController Home) Index(ctx * gin.Context){
		db, err := Connect();
		if err == nil{
			var products []Model.Product
			db.Select(&products,"SELECT * FROM products")
			ctx.HTML(http.StatusOK,"index.html",gin.H{
				"items":products,
			})
		}else{
			fmt.Println(err.Error())
		}
		db.Close()
}
