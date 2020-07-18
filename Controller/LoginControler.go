package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/sessions"
)

type LoginCon struct{

}

var LoginController LoginCon

func (loginCon LoginCon) View (ctx * gin.Context){
	ctx.HTML(http.StatusOK,"login.html",gin.H{
		"error":"",
		"username":"",
		"password":"",
	})
}

func (loginCon LoginCon)Login(ctx * gin.Context){
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if(username != "admin" || password != "adminbttr"){
		ctx.HTML(http.StatusOK,"login.html",gin.H{
			"error":"Wrong username or password",
			"username":username,
			"password":password,
		})
	}else{
		session:= sessions.Default(ctx)
		session.Set("user","admin")
		session.Save()
		ctx.Redirect(http.StatusSeeOther,"/")
		ctx.Abort()
	}
}

func(loginCon LoginCon) Logout(ctx * gin.Context){
	session := sessions.Default(ctx)
	session.Set("user","")
	session.Save()
	ctx.Redirect(http.StatusSeeOther,"/")
}

func(loginCon LoginCon) User(ctx *gin.Context){
	session:= sessions.Default(ctx)
	ctx.JSON(http.StatusOK,gin.H{
		"user":session.Get("user"),
	})
}