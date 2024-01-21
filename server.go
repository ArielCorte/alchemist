package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.File("/", "views/index.html")

	ingredients := []Ingredient{
		{"air", true},
		{"earth", true},
	}

	e.GET("/unlocked_ingredients", func(c echo.Context) error {

		html := ""

		for _, ingredient := range ingredients {
			if ingredient.unlocked {
				html += "<li>" + ingredient.name + "</li>"
			}
		}

		return c.HTML(200, html)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

type Ingredient struct {
	name     string
	unlocked bool
}
