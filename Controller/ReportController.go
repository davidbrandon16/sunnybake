package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReportCon struct{

}

var ReportController ReportCon

func (reportCon ReportCon) View(ctx *gin.Context){
	ctx.HTML(http.StatusOK, "report.html",gin.H{

	})
}
