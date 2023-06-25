package main

import (
	"net/http"

	"github.com/hopk8412/table-recipes-api/configs"

	"github.com/hopk8412/table-recipes-api/routes"

	"github.com/hopk8412/table-recipes-api/controllers"

	"github.com/gin-gonic/gin"

)

// recipe represents data about a recipe.
type recipe struct {
	ID           string   `json:"_id"`
	Title        string   `json:"title"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	AuthorId     string   `json:"authorId"`
}

// recipes slice to seed recipe data.
var recipes = []recipe{
	{ID: "2349028aklsdf", Title: "Braised Beef", Ingredients: []string{"1 Cup Milk", "1 Liter of Cola"}, Instructions: []string{"Drink Milk", "Drink Cola"}, AuthorId: "20394lksdfl"},
	{ID: "asdfjdfkl2323423", Title: "Butter Chicken", Ingredients: []string{"1 Naan Bread", "Steamy Rice"}, Instructions: []string{"Eat Bread", "Eat Rice"}, AuthorId: "dldldl202020"},
	{ID: "ffffff12345", Title: "New Recipe Title", Ingredients: []string{"Ing1", "Ing2"}, Instructions: []string{"Ins1", "Ins2"}, AuthorId: "vvvvvv0293498"},
}

// getRecipes responds with the list of all recipes as JSON.
// func getRecipes(c *gin.Context) {
// 	client := getMongoClient()
// 	recipeCollection := client.Database("table").Collection("recipes")
// 	// recipeTitle := "butter chicken"
// 	cursor, err := recipeCollection.Find(context.TODO(), bson.D{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var results []recipe
// 	if err = cursor.All(context.TODO(), &results); err != nil {
// 		log.Fatal(err)
// 	}
// 	c.IndentedJSON(http.StatusOK, results)
// }

// postRecipes adds a recipe from JSON received in the request body.
func postRecipes(c *gin.Context) {
	var newRecipe recipe

	// Call BindJSON to bind the received JSON to
	// newRecipe.
	if err := c.BindJSON(&newRecipe); err != nil {
		return
	}

	// Add the new recipe to the slice.
	recipes = append(recipes, newRecipe)
	c.IndentedJSON(http.StatusCreated, newRecipe)
}

// getRecipeById locates the recipe whose ID value matches the id
// parameter sent by the client, then returns that recipe as a response.
func getRecipeByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of recipes, looking for
	// a recipe whose ID value matches the parameter.
	for _, a := range recipes {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "recipe not found"})
}

func main() {
	router := gin.Default()

	configs.ConnectDB()

	prefix := "/api/v1"
	routes.RecipeRoutes(router)
	router.GET(prefix+"/recipes", controllers.GetAllRecipes())
	router.GET(prefix+"/recipes/:id", getRecipeByID)
	router.POST(prefix+"/recipes", postRecipes)
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "We couldn't find the page you requested!"})
	})

	router.Run("localhost:8080")
}
