package main

import (
	"github.com/gin-gonic/gin"
	"sunnybake/Controller"
	"github.com/gin-contrib/sessions"
	"net/http"
	"fmt"
	"os"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context){
		session := sessions.Default(ctx)
		user := session.Get("user")
		fmt.Println(user)
		if(user != "admin"){
			ctx.Redirect(http.StatusSeeOther,"/")
			return;
		}else {
			fmt.Println(user)
			ctx.Next();
		}

	}
}

func LoginMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context){
		session := sessions.Default(ctx)
		user := session.Get("user")
		if(user == "admin"){
			ctx.Redirect(http.StatusSeeOther,"/")
			return;
		}else{
			fmt.Println(user)
			ctx.Next();
		}

	}
}

func main() {

	router := gin.Default()
	router.LoadHTMLFiles(
		"View/index.html",
		"View/insert_product.html",
		"View/insert_transaction.html",
		"View/header.html",
		"View/footer.html",
		"View/view_transaction.html",
		"View/change_payment.html",
		"View/order_view.html",
		"View/order_view_detail.html",
		"View/view_receipt.html",
		"View/print_receipt.html",
		"View/report.html",
		"View/login.html",
		"View/update_product.html",
		"View/insert_discount.html",
		"View/view_discount.html",
		"View/update_discount.html",
		"View/update_transaction.html",
	)
	store:= sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("authUser",store))
	router.Static("/asset","./static")
	router.GET("/",Controller.HomeController.Index)
	router.GET("/logout", Controller.LoginController.Logout)
	router.GET("/user",Controller.LoginController.User)
	login := router.Group("login")
	login.Use(LoginMiddleware())
	{
		login.GET("/",Controller.LoginController.View)
		login.POST("/",Controller.LoginController.Login)
	}

	product := router.Group("/product")
	product.Use(AuthMiddleware())
	{
		product.GET("/",Controller.ProductController.Index)
		product.GET("/insert",Controller.ProductController.Index)
		product.POST("/insert",Controller.ProductController.Insert)
		product.GET("/update/:id",Controller.ProductController.ViewUpdate)
		product.POST("/update/:id",Controller.ProductController.Update)
		product.GET("/delete/:id",Controller.ProductController.Delete)
	}
	transaction := router.Group("/transaction")
	transaction.Use(AuthMiddleware())
	{
		transaction.GET("/insert",Controller.TransactionController.Index)
		transaction.POST("/insert",Controller.TransactionController.Insert)
		transaction.GET("/view",Controller.TransactionController.View)
		transaction.GET("/changePayment/:id",Controller.TransactionController.ChangePayment)
		transaction.POST("/changePayment/:id",Controller.TransactionController.InsertPayment)
		transaction.GET("/delete/:id",Controller.TransactionController.Delete)
		transaction.GET("/update/:id",Controller.TransactionController.Update)
		transaction.POST("/update/:id",Controller.TransactionController.UpdateData)
	}

	order := router.Group("/order")
	order.Use(AuthMiddleware())
	{
		order.GET("/view",Controller.OrderController.View)
		order.GET("/viewDetail/:id",Controller.OrderController.ViewDetail)
		order.POST("/viewDetail/:id",Controller.OrderController.ChangeStatus)
	}

	receipt := router.Group("/receipt")
	receipt.Use(AuthMiddleware())
	{
		receipt.GET("/view",Controller.ReceiptController.View)
		receipt.GET("/view/:page",Controller.ReceiptController.ViewPage)
		receipt.GET("/print/:id",Controller.ReceiptController.Print)
	}

	discount := router.Group("/discount")
	discount.Use(AuthMiddleware())
	{
		discount.GET("/view",Controller.DiscountController.View)
		discount.GET("/insert",Controller.DiscountController.Index)
		discount.POST("/insert",Controller.DiscountController.Insert)
		discount.GET("/update/:id",Controller.DiscountController.Update)
		discount.POST("/update/:id",Controller.DiscountController.UpdateData)
		discount.GET("/delete/:id",Controller.DiscountController.Delete)
	}

	report := router.Group("/report")
	report.Use(AuthMiddleware())
	{
		report.GET("/",Controller.ReportController.View)
	}

	port := os.Getenv("PORT")
	router.Run(":"+port)
}
