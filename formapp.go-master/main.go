package main

import (
	"fmt"
	"net/http"

    "github.com/gin-gonic/gin"

    "formapp.go/service"
)

// config
const port = 8000

func main() {
    // initialize Gin engine
    engine := gin.Default()
    engine.LoadHTMLGlob("templates/*.html")

    // routing
    engine.GET("/", rootHandler)
    engine.GET("/name-form", service.NameFormHandler)
    engine.POST("/register-name", service.RegisterNameHandler)

    engine.GET("/start-stateless", service.StartStatelessHandler)
    engine.POST("/name-form-stateless", service.NameFormStatelessHandler)
    engine.POST("/date-form-stateless", service.DateFormStatelessHandler)
    engine.POST("/message-form-stateless", service.MessageFormStatelessHandler)
    engine.POST("/confirmation-stateless", service.ConfirmationStatelessHandler)
    engine.POST("/start-stateless", service.StartStatelessHandler)

    engine.GET("/start-session", service.StartSessionHandler)
    engine.POST("/name-form-session", service.NameFormSessionHandler)
    engine.POST("/date-form-session", service.DateFormSessionHandler)
    engine.POST("/message-form-session", service.MessageFormSessionHandler)
    engine.POST("/confirmation-session", service.ConfirmationSessionHandler)
    engine.POST("/start-session", service.StartSessionHandler)

    engine.GET("/confirmation-a-2-3", service.ConfirmationA_2_3Handler)

    // start server
    engine.Run(fmt.Sprintf(":%d", port))
}

func rootHandler(ctx *gin.Context) {
    // ctx.String(http.StatusOK, "Hello world.")
    ctx.HTML(http.StatusOK, "hello.html", nil)
}