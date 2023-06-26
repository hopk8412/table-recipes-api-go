package models

type MongoUser struct {
	Id string `bson:"_id,omitempty" json:"id,omitempty"`
	FavoriteRecipes []string `json:"favoriteRecipes"`
}