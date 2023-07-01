package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hopk8412/table-recipes-api/configs"
	"github.com/hopk8412/table-recipes-api/models"
	"github.com/hopk8412/table-recipes-api/responses"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var recipeCollection *mongo.Collection = configs.GetCollection(configs.DB, "recipes")
var usersCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

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

func GetRecipeById() gin.HandlerFunc {
	return func(c *gin.Context) {
		recipeId := c.Param("id")
		log.Println("Attempting to retrieve recipe with ID: ", recipeId)		
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var recipe models.Recipe
		defer cancel()

		result := recipeCollection.FindOne(ctx, bson.M{"_id": c.Param("id")})
		
		err := result.Decode(&recipe) 
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		} 
		
		c.JSON(http.StatusOK, responses.RecipeResponse{Status: http.StatusOK, Message: "Successfully fetched recipe with ID " + c.Param("id"), Data: map[string]interface{}{"data": recipe}})
		
	}
}

func GetRecipesByAuthorId() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeaderValue := c.Request.Header["Authorization"][0]
		// Need to use auth header value and pass it on to Keycloak to validate user
		log.Println("Validating user before fetching list of recipes created by user...")
		keycloakHttpClient := &http.Client{}
		req, err := http.NewRequest("GET", configs.EnvUserInfoURI(), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error validating user", Data: map[string]interface{}{"data": err.Error()}})
		}
		req.Header.Add("Authorization", authHeaderValue)
		resp, err := keycloakHttpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error validating user", Data: map[string]interface{}{"data": err.Error()}})
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error validating user", Data: map[string]interface{}{"data": err.Error()}})
		}
		
		var keycloakUser models.KeycloakUser
		json.Unmarshal(bodyBytes, &keycloakUser)

		// User should be validated at this point - fetch recipes by provided author's ID
		log.Println("Attempting to retrieve recipes created by user with ID: ", keycloakUser.Sub)		
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var recipes []models.Recipe
		defer cancel()

		results, err := recipeCollection.Find(ctx, bson.M{"authorId": keycloakUser.Sub})
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
		c.JSON(http.StatusOK, responses.RecipeResponse{Status: http.StatusOK, Message: "Successfully fetched all recipes created by user with ID: " + keycloakUser.Sub, Data: map[string]interface{}{"data": recipes}})
	}
}

func PostRecipe() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var recipe models.Recipe
		defer cancel()

		//validate request body
		if err := c.BindJSON(&recipe); err != nil {
			c.JSON(http.StatusBadRequest, responses.RecipeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//TODO: validate required fields

		newRecipe := models.Recipe{
			Id: primitive.NewObjectID().Hex(),
			Title: recipe.Title,
			Ingredients: recipe.Ingredients,
			Instructions: recipe.Instructions,
			AuthorId: recipe.AuthorId,
			ImageLinks: recipe.ImageLinks,
		}
		result, err := recipeCollection.InsertOne(ctx, newRecipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, responses.RecipeResponse{Status: http.StatusCreated, Message: "Successfully created recipe!", Data: map[string]interface{}{"data": result}})
	}
}

func DeleteRecipeById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		recipeId := c.Param("id")
		defer cancel()

		result, err := recipeCollection.DeleteOne(ctx, bson.M{"_id": recipeId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.RecipeResponse{Status: http.StatusOK, Message: "Successfully deleted recipe with ID " + recipeId, Data: map[string]interface{}{"data": result}})

	}
}

func SearchForRecipes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var recipes []models.Recipe
		var searchQuery models.SearchQuery
		defer cancel()

		//validate request body
		if err := c.BindJSON(&searchQuery); err != nil {
			c.JSON(http.StatusBadRequest, responses.RecipeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		results, err := recipeCollection.Find(ctx, bson.D{{Key: "title", Value: primitive.Regex{Pattern: searchQuery.SearchTerm, Options: "i"}}})
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
		c.JSON(http.StatusOK, responses.RecipeResponse{Status: http.StatusOK, Message: "Successfully fetched all recipes with title containing '" + searchQuery.SearchTerm + "'!", Data: map[string]interface{}{"data": recipes}})
	}
}

func AddOrRemoveRecipeToUserFavorites() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate user token before continuing here...
		authHeaderValue := c.Request.Header["Authorization"][0]
		// Need to use auth header value and pass it on to Keycloak to validate user
		log.Println("Validating user before fetching list of recipes created by user...")
		keycloakHttpClient := &http.Client{}
		req, err := http.NewRequest("GET", configs.EnvUserInfoURI(), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error validating user", Data: map[string]interface{}{"data": err.Error()}})
		}
		req.Header.Add("Authorization", authHeaderValue)
		resp, err := keycloakHttpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error validating user", Data: map[string]interface{}{"data": err.Error()}})
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error validating user", Data: map[string]interface{}{"data": err.Error()}})
		}
		
		var keycloakUser models.KeycloakUser
		json.Unmarshal(bodyBytes, &keycloakUser)

		// At this point, user is validated, continue with adding/removing favorite for user
		if keycloakUser.Sub == c.Param("id") {
		// Perform check on Token sub value matching provided userId from request...
		if c.Param("id") == keycloakUser.Sub {
			log.Println("Token was validated for user with ID ", keycloakUser.Sub, " - persisting user to mongo if not already there...")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			var userRecipeOperation models.UserRecipeOperation
			var mongoUser models.MongoUser
			defer cancel()

			//validate request body
			if err := c.BindJSON(&userRecipeOperation); err != nil {
				c.JSON(http.StatusBadRequest, responses.RecipeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}

			err := usersCollection.FindOne(ctx, bson.M{"_id": c.Param("id")}).Decode(&mongoUser)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					log.Println("User does not exist, insert now and add provided recipe to their favorites...")
					newMongoUser := models.MongoUser{
						Id: c.Param("id"),
						FavoriteRecipes: []string{userRecipeOperation.RecipeId},
					}
					// Insert new user record
					result, err := usersCollection.InsertOne(ctx, newMongoUser)
					if err != nil {
						c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
						return
					}
					c.JSON(http.StatusCreated, responses.RecipeResponse{Status: http.StatusCreated, Message: "Successfully created recipe!", Data: map[string]interface{}{"data": result}})
					return
				} else {
					c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					return
				}
			}
			// At this point, we know the mongo user exists, perform requested operation and save change
			log.Println("User already exists, adding/removing recipe now...")
			if userRecipeOperation.IsAddingFavorite {
				mongoUser.FavoriteRecipes = append(mongoUser.FavoriteRecipes, userRecipeOperation.RecipeId)
				update := bson.M{"favoriteRecipes": mongoUser.FavoriteRecipes}
				result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": c.Param("id")}, bson.M{"$set": update})
				if err != nil {
					c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					return
				}
				c.JSON(http.StatusOK, responses.RecipeResponse{Status: http.StatusOK, Message: "Successfully added recipe to user favorites!", Data: map[string]interface{}{"data": result}})
			} else {
				log.Println("Removing recipe with ID ", userRecipeOperation.RecipeId, " from users favorites...")
				for i, v := range mongoUser.FavoriteRecipes {
					if v == userRecipeOperation.RecipeId {
						mongoUser.FavoriteRecipes = append(mongoUser.FavoriteRecipes[:i], mongoUser.FavoriteRecipes[i+1:]...)
						break
					}
				}
				update := bson.M{"favoriteRecipes": mongoUser.FavoriteRecipes}
				result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": c.Param("id")}, bson.M{"$set": update})
				if err != nil {
					c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					return
				}
				c.JSON(http.StatusOK, responses.RecipeResponse{Status: http.StatusOK, Message: "Successfully removed recipe from user favorites!", Data: map[string]interface{}{"data": result}})
			}
		}
	}
	}
}

func GetUserFavoriteRecipes() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate user token before continuing here...
		authHeaderValue := c.Request.Header["Authorization"][0]
		// Need to use auth header value and pass it on to Keycloak to validate user
		log.Println("Validating user before fetching list of recipes created by user...")
		keycloakHttpClient := &http.Client{}
		req, err := http.NewRequest("GET", configs.EnvUserInfoURI(), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error validating user", Data: map[string]interface{}{"data": err.Error()}})
		}
		req.Header.Add("Authorization", authHeaderValue)
		resp, err := keycloakHttpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error validating user", Data: map[string]interface{}{"data": err.Error()}})
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error validating user", Data: map[string]interface{}{"data": err.Error()}})
		}
		
		var keycloakUser models.KeycloakUser
		json.Unmarshal(bodyBytes, &keycloakUser)

		// At this point, user is validated, continue with adding/removing favorite for user
		log.Println("User was successfully validated, getting users favorited recipes...")
		if keycloakUser.Sub == c.Param("id") {
			// Query users collection by user ID to get favorite recipes
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			var mongoUser models.MongoUser
			defer cancel()

			result := usersCollection.FindOne(ctx, bson.M{"_id": c.Param("id")})
		
			err := result.Decode(&mongoUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error validating user", Data: map[string]interface{}{"data": err.Error()}})
			}
			// We now have the mongo user, so use their favorite recipe ID slice to query recipe collection...
			filter := bson.M{"_id": bson.M{"$in": mongoUser.FavoriteRecipes}}
			var recipes []models.Recipe
			defer cancel()

			favoritedRecipes, err := recipeCollection.Find(ctx, filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			//read from mongo optimally
			defer favoritedRecipes.Close(ctx)
			for favoritedRecipes.Next(ctx) {
				var singleRecipe models.Recipe
				if err = favoritedRecipes.Decode(&singleRecipe); err != nil {
					c.JSON(http.StatusInternalServerError, responses.RecipeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				}
				recipes = append(recipes, singleRecipe)
				
			}
			c.JSON(http.StatusOK, responses.RecipeResponse{Status: http.StatusOK, Message: "Successfully fetched all recipes!", Data: map[string]interface{}{"data": recipes}})
		}
	}
}