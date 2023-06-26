package models

type UserRecipeOperation struct {
	RecipeId string `json:"recipeId"`
	IsAddingFavorite bool `json:"isAddingFavorite"`
}