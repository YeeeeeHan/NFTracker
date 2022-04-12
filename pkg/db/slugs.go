package db

import (
	"github.com/go-pg/pg/v10"
)

type Slug struct {
	ID               int64  `json:"id"`
	SlugName         int64  `json:"slug_name"`
	FloorPrice       string `json:"floor_price"`
	OneDayAvgPrice   string `json:"one_day_average_price"`
	SevenDayAvgPrice int64  `json:"seven_day_average_price"`
}

func CreateSlug(db *pg.DB, req *Slug) (*Slug, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return nil, err
	}

	slug := &Slug{}
	err = db.Model(slug).
		Relation("Users").
		Where("slug.id = ?", req.ID).
		Select()

	return slug, err
}

func GetSlug(db *pg.DB, homeID string) (*Slug, error) {
	slug := &Slug{}
	err := db.Model(slug).
		Relation("Users").
		Where("slug.id = ?", homeID).
		Select()

	return slug, err
}

func GetSlugs(db *pg.DB) ([]*Slug, error) {
	slugs := make([]*Slug, 0)
	err := db.Model(&slugs).
		Relation("Users").
		Select()

	return slugs, err
}

func UpdateSlug(db *pg.DB, s *Slug) (*Slug, error) {
	_, err := db.Model(s).WherePK().Update()
	if err != nil {
		return nil, err
	}

	slug := &Slug{}
	err = db.Model(slug).
		Relation("Users").
		Where("slug.id = ?", s.ID).
		Select()

	return slug, err
}

func DeleteSlug(db *pg.DB, slugID int64) error {
	slug := &Slug{
		ID: slugID,
	}

	err := db.Model(slug).
		Relation("Users").
		Where("slug.id = ?", slug.ID).
		Select()
	if err != nil {
		return err
	}

	_, err = db.Model(slug).WherePK().Delete()

	return err
}
