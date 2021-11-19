package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")

	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/", service.Home)
	engine.POST("/", service.Home)
	engine.GET("/list", service.TaskList)
	engine.GET("/task/:id", service.ShowTask) // ":id" is a parameter

	engine.GET("/create", service.Create)
	engine.POST("/create/confirm", service.CreateConfirm)
	engine.GET("/task/:id/edit", service.Edit)
	engine.POST("/task/:id/edit/confirm", service.EditConfirm)
	engine.POST("/task/:id/delete", service.Delete)
	engine.POST("/task/:id/share", service.Share)
	engine.GET("/search", service.Search)

	engine.GET("/login", service.LogIn)
	engine.POST("/login/confirm", service.LogInConfirm)
	engine.GET("/edit-account", service.EditAccount)
	engine.POST("/edit-account/confirm", service.EditAccountConfirm)
	engine.GET("/signup", service.SignUp)
	engine.POST("/signup/confirm", service.SignUpConfirm)
	engine.GET("/delete-account", service.DeleteAccount)
	engine.POST("/delete-account/confirm", service.DeleteAccountConfirm)

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
