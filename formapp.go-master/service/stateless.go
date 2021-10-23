package service

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func StartStatelessHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "start.stateless.html", nil)
}

func NameFormStatelessHandler(ctx *gin.Context) {
	var data FormData
	ctx.Bind(&data)
	ctx.HTML(http.StatusOK, "name_form.stateless.html", &data)
}

func DateFormStatelessHandler(ctx *gin.Context) {
	var data FormData
	ctx.Bind(&data)
	ctx.HTML(http.StatusOK, "date_form.stateless.html", &data)
}

func MessageFormStatelessHandler(ctx *gin.Context) {
	var data FormData
	ctx.Bind(&data)
	ctx.HTML(http.StatusOK, "message_form.stateless.html", &data)
}

func ConfirmationStatelessHandler(ctx *gin.Context) {
	var data FormData
	ctx.Bind(&data)
	ctx.HTML(http.StatusOK, "confirmation.stateless.html", &data)
}