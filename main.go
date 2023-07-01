package main

import (
	"net/http"

	"github.com/hopk8412/table-recipes-api/configs"
	"golang.org/x/exp/slices"

	"github.com/hopk8412/table-recipes-api/routes"

	"github.com/hopk8412/table-recipes-api/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	configs.ConnectDB()

	router.Use(corsMiddleware())

	prefix := "/api/v1"
	routes.RecipeRoutes(router)
	router.GET(prefix+"/recipes", controllers.GetAllRecipes())
	router.GET(prefix+"/recipes/:id", controllers.GetRecipeById())
	router.GET(prefix+"/recipes/me", controllers.GetRecipesByAuthorId())
	router.GET(prefix+"/users/:id/recipes", controllers.GetUserFavoriteRecipes())
	router.POST(prefix+"/recipes", controllers.PostRecipe())
	router.POST(prefix+"/recipes/search", controllers.SearchForRecipes())
	router.POST(prefix+"/users/:id/recipes", controllers.AddOrRemoveRecipeToUserFavorites())
	router.DELETE(prefix+"/recipes/:id", controllers.DeleteRecipeById())
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "We couldn't find the page you requested!"})
	})

	router.Run(":8080")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if origin := c.Request.Header.Get("Origin"); slices.Contains(configs.AllowedOrigins(), origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}
		c.Next()
	}
}