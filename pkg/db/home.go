package db

type Home struct {
	ID          int64  `json:"id"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Agent       *Agent `pg:"rel:has-one" json:"agent"`
}
