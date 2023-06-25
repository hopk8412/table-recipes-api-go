package models

type Recipe struct {
	Id           string   `bson:"_id,omitempty" json:"id,omitempty"`
	Title        string   `json:"title,omitempty"`
	Ingredients  []string `json:"ingredients,omitempty"`
	Instructions []string `json:"instructions,omitempty"`
	AuthorId     string   `json:"authorId,omitempty"`
}
