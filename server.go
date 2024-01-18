package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./alchemist.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e := echo.New()
	e.File("/", "views/index.html")
	e.Logger.Fatal(e.Start(":1323"))
}
