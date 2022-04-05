package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sjcl/shionstagram-backend/web/model"
	"github.com/sjcl/shionstagram-backend/web/handler"
)

func main() {
	db := sqlx.MustConnect("mysql", "shion:password@tcp(db:3306)/shionstagram_db?parseTime=true")
	defer db.Close()

	m := model.NewModel(db)
	h := handler.NewHandler(m)
	
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/message", h.PostMessage)

	e.Logger.Fatal(e.Start(":8083"))
}