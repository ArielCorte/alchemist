package main

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.File("/", "views/index.html")

	json_Ingredients, err := os.ReadFile("resources/map_ingredients.json")
	if err != nil {
		panic(err)
	}

	var ingredients map[int]Ingredient

	err = json.Unmarshal(json_Ingredients, &ingredients)
	if err != nil {
		panic(err)
	}

	unlocked_ingredients := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 345, 235, 6, 456, 123}

	soup_ing := []int{}
	e.Static("/css", "css")

	e.GET("/unlocked_ingredients", func(c echo.Context) error {
		search_query := c.QueryParam("search")

		curr_ing := make(map[int]Ingredient)

		if search_query != "" {
			for _, id := range unlocked_ingredients {
				if strings.HasPrefix(strings.ToLower(ingredients[id].Name), strings.ToLower(search_query)) {
					curr_ing[id] = ingredients[id]
				}
			}
		} else {
			for _, id := range unlocked_ingredients {
				curr_ing[id] = ingredients[id]
			}
		}

		html := ""

		sorted_ingredients := make(PairList, 0, len(curr_ing))

		for id, ingredient := range curr_ing {
			sorted_ingredients = append(sorted_ingredients, Pair{id, ingredient.Name})
		}
		sort.Sort(sorted_ingredients)

		first := true

		for _, i := range sorted_ingredients {
			trigger := "click"
			ingredient := curr_ing[i.Key]
			if containi(unlocked_ingredients, ingredient.Id) {
				if first {
					trigger += ", keyup[keyCode==13] from:#search"
					first = false
				}
				html += "<li hx-target='#current-ingredients' hx-trigger='" + trigger + "' hx-post='/add_ingredient?ingredient=" + strconv.Itoa(ingredient.Id) + "'>" + ingredient.Name + "</li>"

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
			return c.HTML(200, "<div>"+ingredients[soup_ing[0]].Name+"</div><div hx-target='#result-soup' hx-get='get_result' hx-trigger='load'>"+ingredients[soup_ing[1]].Name+"</div>")
		}

		return c.HTML(200, "<div hx-target='#result-soup, input' hx-put='clear_result' hx-trigger='load'>"+ingredients[id].Name+"</div><div></div>")
	})

	e.GET("/get_result", func(c echo.Context) error {
		crafted_ing, err := craftSoup(ingredients, soup_ing)
		soup_ing = []int{}
		if err != nil {
			return c.String(200, err.Error())
		}
		unlocked_ingredients = append(unlocked_ingredients, crafted_ing.Id)
		return c.HTML(200, "<span hx-get='unlocked_ingredients' hx-target='#unlocked-ingredients' hx-trigger='load'>"+crafted_ing.Name+"</span>")
	})

	e.PUT("/clear_result", func(c echo.Context) error {
		return c.String(200, "")
	})

	e.DELETE("/reset", func(c echo.Context) error {
		soup_ing = []int{}
		return c.HTML(200, "<div hx-target='#result-soup' hx-put='clear_result' hx-trigger='load'></div><div></div>")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func craftSoup(ings map[int]Ingredient, soup_ing []int) (Ingredient, error) {
	if len(soup_ing) < 2 {
		return Ingredient{}, errors.New("not enough ingredients")
	}
	// check if ingredients are valid
	for _, ingredient := range ings {
		for _, recipe := range ingredient.Parents {
			if soup_ing[0] != soup_ing[1] {
				if containi(recipe, soup_ing[0]) && containi(recipe, soup_ing[1]) {
					return ingredient, nil
				}
			} else {
				if len(recipe) == 2 && recipe[0] == soup_ing[0] && recipe[1] == soup_ing[1] {
					return ingredient, nil
				}
			}
		}
	}
	// if valid, craft soup
	return Ingredient{}, errors.New("no ingredient found")
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
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Parents [][]int  `json:"parents"`
	Tags    []string `json:"tags"`
}

type Pair struct {
	Key   int
	Value string
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
