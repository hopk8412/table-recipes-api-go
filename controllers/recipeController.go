package controllers

import (
	"context"
	"github.com/hopk8412/table-recipes-api/configs"
	"github.com/hopk8412/table-recipes-api/models"
	"github.com/hopk8412/table-recipes-api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

var recipeCollection *mongo.Collection = configs.GetCollection(configs.DB, "recipes")

func GetAllRecipes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var recipes []models.Recipe
		defer cancel()

		results, err := recipeCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		//read from mongo optimally
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleRecipe models.Recipe
			if err = results.Decode(&singleRecipe); err != nil {
				c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}
			recipes = append(recipes, singleRecipe)
		}
		c.JSON(http.StatusOK, responses.RecipeResponse{Status: http.StatusOK, Message: "Successfully fetched all recipes!", Data: map[string]interface{}{"data": recipes}})
	}
}
