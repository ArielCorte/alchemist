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
		{"fire", true},
		{"water", true},
		{"alcohol", false},
	}

	soup_ing := []Ingredient{}

	e.Static("/css", "css")

	e.GET("/unlocked_ingredients", func(c echo.Context) error {

		html := ""

		for _, ingredient := range ingredients {
			if ingredient.unlocked {
				html += "<li hx-swap='outerHTML' hx-target='#available-soup' hx-post='/add_ingredient?ingredient=" + ingredient.name + "'>" + ingredient.name + "</li>"
			}
		}

		return c.HTML(200, html)
	})

	e.POST("/add_ingredient", func(c echo.Context) error {
		ingredient := c.QueryParam("ingredient")

		soup_ing = append(soup_ing, Ingredient{ingredient, true})

		if len(soup_ing) == 2 {
			// send request to craft soup
			return c.HTML(200, "<li hx-target='#result-soup' hx-get='get_result' hx-trigger='load delay:1s'>"+ingredient+"</li>")
		}

		return c.HTML(200, "<li>"+ingredient+"</li><div id='available-soup'></div>")
	})

	e.GET("/get_result", func(c echo.Context) error {
		crafted_ing, err := craftSoup(soup_ing)
		soup_ing = []Ingredient{}
		if err != nil {
			return c.String(200, "no ingredient found")
		}
		return c.HTML(200, "<li>"+crafted_ing.name+"</li>")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func craftSoup(ingredients []Ingredient) (Ingredient, error) {
	// check if ingredients are valid
	// if valid, craft soup
	return Ingredient{"alcohol", true}, nil
}

type Ingredient struct {
	name     string
	unlocked bool
}
