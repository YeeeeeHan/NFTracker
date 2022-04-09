package db

import (
	"github.com/go-pg/pg/v10"
)

type User struct {
	ID       int64 `json:"id"`
	Username int64 `json:"username"`
}

func CreateUser(db *pg.DB, req *User) (*User, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = db.Model(user).
		Relation("User").
		Where("user.id = ?", req.ID).
		Select()

	return user, err
}

func GetUser(db *pg.DB, homeID string) (*User, error) {
	user := &User{}
	err := db.Model(user).
		Relation("User").
		Where("user.id = ?", homeID).
		Select()

	return user, err
}

func GetUsers(db *pg.DB) ([]*User, error) {
	slugs := make([]*User, 0)
	err := db.Model(&slugs).
		Relation("User").
		Select()

	return slugs, err
}

func UpdateUser(db *pg.DB, u *User) (*User, error) {
	_, err := db.Model(u).WherePK().Update()
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = db.Model(user).
		Relation("User").
		Where("user.id = ?", u.ID).
		Select()

	return user, err
}

func DeleteUser(db *pg.DB, userID int64) error {
	user := &User{
		ID: userID,
	}

	err := db.Model(user).
		Relation("User").
		Where("user.id = ?", user.ID).
		Select()
	if err != nil {
		return err
	}

	_, err = db.Model(user).WherePK().Delete()

	return err
}
