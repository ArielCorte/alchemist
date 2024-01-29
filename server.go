package main

import (
	"errors"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.File("/", "views/index.html")

	ingredients := map[int]Ingredient{
		1: {1, "air", [][]int{}, []string{}},
		2: {2, "earth", [][]int{}, []string{}},
		3: {3, "fire", [][]int{}, []string{}},
		4: {4, "water", [][]int{}, []string{}},
		5: {5, "rain", [][]int{{1, 4}, {8, 4}}, []string{}},
		6: {6, "pressure", [][]int{{1, 1}}, []string{}},
		7: {7, "steam", [][]int{{4, 3}, {9, 4}}, []string{}},
		8: {8, "cloud", [][]int{{1, 7}}, []string{}},
	}

	unlocked_ingredients := []int{}

	soup_ing := []int{}
	e.Static("/css", "css")

	e.GET("/unlocked_ingredients", func(c echo.Context) error {

		html := ""

		sorted_ingredients := make(PairList, 0, len(ingredients))

		for id, ingredient := range ingredients {
			sorted_ingredients = append(sorted_ingredients, Pair{id, ingredient.Name})
		}
		sort.Sort(sorted_ingredients)

		for _, i := range sorted_ingredients {
			ingredient := ingredients[i.Key]
			if containi(unlocked_ingredients, ingredient.Id) {
				html += "<li hx-target='#current-ingredients' hx-post='/add_ingredient?ingredient=" + ingredient.Name + "'>" + ingredient.Name + "</li>"
			}
		}

		return c.HTML(200, html)
	})

	e.POST("/add_ingredient", func(c echo.Context) error {
		string_id := c.QueryParam("ingredient")
		//convert string to int
		id, err := strconv.Atoi(string_id)
		if err != nil {
			return c.String(200, err.Error())
		}

		soup_ing = append(soup_ing, id)

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
		unlocked_ingredients = append(unlocked_ingredients, crafted_ing.Id)
		ingredients[crafted_ing_name] = crafted_ing
		return c.HTML(200, "<span hx-get='unlocked_ingredients' hx-target='#unlocked-ingredients' hx-trigger='load'>"+crafted_ing_name+"</span>")
	})

	e.PUT("/clear_result", func(c echo.Context) error {
		return c.String(200, "")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func craftSoup(ings map[string]Ingredient, soup_ing []int) (string, error) {
	if len(soup_ing) < 2 {
		return "", errors.New("not enough ingredients")
	}
	// check if ingredients are valid
	for _, ingredient := range ings {
		for _, recipe := range ingredient.Parents {
			if soup_ing[0] != soup_ing[1] {
				if contains(recipe, soup_ing[0]) && contains(recipe, soup_ing[1]) {
					return ingredient.Name, nil
				}
			} else {
				if len(recipe) == 2 && recipe[0] == soup_ing[0] && recipe[1] == soup_ing[1] {
					return ingredient.Name, nil
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

func containi(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type Ingredient struct {
	Id      int
	Name    string
	Parents [][]int
	Tags    []string
}

type Pair struct {
	Key   int
	Value string
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
