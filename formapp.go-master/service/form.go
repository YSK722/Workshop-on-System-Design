package service

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type FormData struct {
	Name string `form:"name"`
	Birthday string `form:"birthday"`
	Message string `form:"message"`
}

func NameFormHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "name_form.html", nil)
}

func RegisterNameHandler(ctx *gin.Context) {
	var data FormData
	ctx.Bind(&data)
	ctx.HTML(http.StatusOK, "result.html", &data)
}

func ConfirmationA_2_3Handler(ctx *gin.Context) {
	var data FormData
	ctx.Bind(&data)
	ctx.HTML(http.StatusOK, "confirmationA_2_3.html", &data)
}