package main

import (
	"errors"
	"sort"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.File("/", "views/index.html")

	ingredients := map[string]Ingredient{
		"air":      {"air", true, [][]string{}},
		"earth":    {"earth", true, [][]string{}},
		"fire":     {"fire", true, [][]string{}},
		"water":    {"water", true, [][]string{}},
		"rain":     {"rain", false, [][]string{{"air", "water"}, {"cloud", "water"}}},
		"pressure": {"pressure", false, [][]string{{"air", "air"}}},
		"steam":    {"steam", false, [][]string{{"water", "fire"}, {"energy", "water"}}},
		"cloud":    {"cloud", false, [][]string{{"air", "steam"}}},
	}

	soup_ing := []string{}

	e.Static("/css", "css")

	e.GET("/unlocked_ingredients", func(c echo.Context) error {

		html := ""

		sorted_ingredients := make([]string, 0, len(ingredients))

		for ingredient := range ingredients {
			sorted_ingredients = append(sorted_ingredients, ingredient)
		}
		sort.Strings(sorted_ingredients)

		for _, i := range sorted_ingredients {
			ingredient := ingredients[i]
			if ingredient.unlocked {
				html += "<li hx-target='#current-ingredients' hx-post='/add_ingredient?ingredient=" + ingredient.name + "'>" + ingredient.name + "</li>"
			}
		}

		return c.HTML(200, html)
	})

	e.POST("/add_ingredient", func(c echo.Context) error {
		ingredient := c.QueryParam("ingredient")

		soup_ing = append(soup_ing, ingredient)

		if len(soup_ing) >= 2 {
			// send request to craft soup
			return c.HTML(200, "<div>"+soup_ing[0]+"</div><div hx-target='#result-soup' hx-get='get_result' hx-trigger='load'>"+soup_ing[1]+"</div>")
		}

		return c.HTML(200, "<div hx-target='#result-soup' hx-put='clear_result' hx-trigger='load'>"+ingredient+"</div><div></div>")
	})

	e.GET("/get_result", func(c echo.Context) error {
		crafted_ing_name, err := craftSoup(ingredients, soup_ing)
		soup_ing = []string{}
		if err != nil {
			return c.String(200, err.Error())
		}
		crafted_ing := ingredients[crafted_ing_name]
		crafted_ing.unlocked = true
		ingredients[crafted_ing_name] = crafted_ing
		return c.HTML(200, "<span hx-get='unlocked_ingredients' hx-target='#unlocked-ingredients' hx-trigger='load'>"+crafted_ing_name+"</span>")
	})

	e.PUT("/clear_result", func(c echo.Context) error {
		return c.String(200, "")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func craftSoup(ings map[string]Ingredient, soup_ing []string) (string, error) {
	if len(soup_ing) < 2 {
		return "", errors.New("not enough ingredients")
	}
	// check if ingredients are valid
	for _, ingredient := range ings {
		for _, recipe := range ingredient.ingredients {
			if soup_ing[0] != soup_ing[1] {
				if contains(recipe, soup_ing[0]) && contains(recipe, soup_ing[1]) {
					return ingredient.name, nil
				}
			} else {
				if len(recipe) == 2 && recipe[0] == soup_ing[0] && recipe[1] == soup_ing[1] {
					return ingredient.name, nil
				}
			}
		}
	}
	// if valid, craft soup
	return "", errors.New("no ingredient found")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type Ingredient struct {
	name        string
	unlocked    bool
	ingredients [][]string
}
