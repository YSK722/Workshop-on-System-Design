package service

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func StartSessionHandler(ctx *gin.Context) {
	uuidObj, _ := uuid.NewRandom()
	ctx.SetCookie("id", uuidObj.String(), 100, "/", "localhost", false, true)
	ctx.HTML(http.StatusOK, "start.session.html", nil)
}

func NameFormSessionHandler(ctx *gin.Context) {
	_, err := ctx.Cookie("id")
	if err != nil {
		ctx.String(http.StatusOK, fmt.Sprintf("Hello, Guest. GO BACK!!"))
	} else {
		ctx.HTML(http.StatusOK, "name_form.session.html", nil)
	}
}

func DateFormSessionHandler(ctx *gin.Context) {
	_, err := ctx.Cookie("id")
	if err != nil {
		ctx.String(http.StatusOK, fmt.Sprintf("Hello, Guest. GO BACK!!"))
	} else {
		name, ok := ctx.GetPostForm("name")
		if ok == true {
			ctx.SetCookie("name", name, 100, "/", "localhost", false, true)
		}
		ctx.HTML(http.StatusOK, "date_form.session.html", nil)
	}
}

func MessageFormSessionHandler(ctx *gin.Context) {
	_, err := ctx.Cookie("id")
	if err != nil {
		ctx.String(http.StatusOK, fmt.Sprintf("Hello, Guest. GO BACK!!"))
	} else {
		birthday, ok := ctx.GetPostForm("birthday")
		if ok == true {
			ctx.SetCookie("birthday", birthday, 100, "/", "localhost", false, true)
		}
		ctx.HTML(http.StatusOK, "message_form.session.html", nil)
	}
}

func ConfirmationSessionHandler(ctx *gin.Context) {
	_, err := ctx.Cookie("id")
	if err != nil {
		ctx.String(http.StatusOK, fmt.Sprintf("Hello, Guest. GO BACK!!"))
	} else {
		message, _ := ctx.GetPostForm("message")
		ctx.SetCookie("message", message, 100, "/", "localhost", false, true)
		name, _ := ctx.Cookie("name")
		birthday, _ := ctx.Cookie("birthday")
		data := FormData{name, birthday, message}
		ctx.HTML(http.StatusOK, "confirmation.session.html", &data)
	}
}