package main

import (
	"github.com/gin-gonic/gin"
	"sunnybake/Controller"
	"github.com/gin-contrib/sessions"
	"net/http"
	"fmt"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context){
		session := sessions.Default(ctx)
		user := session.Get("user")
		fmt.Println(user)
		if(user != "admin"){
			ctx.Redirect(http.StatusSeeOther,"/")
			return;
		}
		fmt.Println(user)
		ctx.Next();

	}
}

func LoginMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context){
		session := sessions.Default(ctx)
		user := session.Get("user")
		if(user == "admin"){
			ctx.Redirect(http.StatusSeeOther,"/")
			return;
		}
		fmt.Println(user)
		ctx.Next();

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
	}
	transaction := router.Group("/transaction")
	{
		transaction.GET("/insert",Controller.TransactionController.Index)
		transaction.POST("/insert",Controller.TransactionController.Insert)
		transaction.GET("/view",Controller.TransactionController.View)
		transaction.GET("/changePayment/:id",Controller.TransactionController.ChangePayment)
		transaction.POST("/changePayment/:id",Controller.TransactionController.InsertPayment)
	}

	order := router.Group("/order")
	{
		order.GET("/view",Controller.OrderController.View)
		order.GET("/viewDetail/:id",Controller.OrderController.ViewDetail)
		order.POST("/viewDetail/:id",Controller.OrderController.ChangeStatus)
	}

	receipt := router.Group("/receipt")
	{
		receipt.GET("/view",Controller.ReceiptController.View)
		receipt.GET("/print/:id",Controller.ReceiptController.Print)
	}

	report := router.Group("/report")
	{
		report.GET("/",Controller.ReportController.View)
	}


	router.Run(":8081")
}
