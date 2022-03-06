package db

import (
	"github.com/go-pg/pg/v10"
)

type Home struct {
	ID          int64  `json:"id"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
	Address     string `json:"address"`
	AgentID     int64  `json:"agent_id"`
	Agent       *Agent `pg:"rel:has-one" json:"agent"`
}

func CreateHome(db *pg.DB, req *Home) (*Home, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return nil, err
	}

	home := &Home{}
	err = db.Model(home).
		Relation("Agent").
		Where("home.id = ?", req.ID).
		Select()

	return home, err
}

func GetHome(db *pg.DB, homeID string) (*Home, error) {
	home := &Home{}
	err := db.Model(home).
		Relation("Agent").
		Where("home.id = ?", homeID).
		Select()

	return home, err
}

func GetHomes(db *pg.DB) ([]*Home, error) {
	homes := make([]*Home, 0)
	err := db.Model(&homes).
		Relation("Agent").
		Select()

	return homes, err
}

func UpdateHome(db *pg.DB, req *Home) (*Home, error) {
	_, err := db.Model(req).WherePK().Update()
	if err != nil {
		return nil, err
	}

	home := &Home{}
	err = db.Model(home).
		Relation("Agent").
		Where("home.id = ?", req.ID).
		Select()

	return home, err
}

func DeleteHome(db *pg.DB, homeID int64) error {
	home := &Home{
		ID: homeID,
	}

	err := db.Model(home).
		Relation("Agent").
		Where("home.id = ?", home.ID).
		Select()
	if err != nil {
		return err
	}

	_, err = db.Model(home).WherePK().Delete()

	return err
}
