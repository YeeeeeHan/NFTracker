package db

import (
	"github.com/go-pg/pg/v10"
)

type Slugs struct {
	SlugName   string  `pg:",pk" json:"slug_name"`
	FloorPrice float64 `json:"floor_price"`
}

func (s *Slugs) String() string {
	return "SlugName=" + s.SlugName
}

func CreateSlug(db *pg.DB, req *Slugs) (*Slugs, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return nil, err
	}

	slug := &Slugs{}
	err = db.Model(slug).
		Where("slugs.slug_name = ?", req.SlugName).
		Select()

	return slug, err
}

func GetSlug(db *pg.DB, slugName string) (*Slugs, error) {
	slug := &Slugs{}
	err := db.Model(slug).
		Where("slugs.slug_name = ?", slugName).
		Select()

	return slug, err
}

func GetAllSlugs(db *pg.DB) ([]*Slugs, error) {
	slugs := make([]*Slugs, 0)
	err := db.Model(&slugs).
		Column("slug_name").
		Select()

	return slugs, err
}

func UpdateSlug(db *pg.DB, s *Slugs) (*Slugs, error) {
	_, err := db.Model(s).WherePK().Update()
	if err != nil {
		return nil, err
	}

	slug := &Slugs{}
	err = db.Model(slug).
		Where("slugs.slug_name = ?", s.SlugName).
		Select()

	return slug, err
}

func DeleteSlug(db *pg.DB, slugID string) error {
	slug := &Slugs{
		SlugName: slugID,
	}

	err := db.Model(slug).
		Relation("Users").
		Where("slug.id = ?", slug.SlugName).
		Select()
	if err != nil {
		return err
	}

	_, err = db.Model(slug).WherePK().Delete()

	return err
}
