package db

import (
	"github.com/go-pg/pg/v10"
)

type Users struct {
	Username string `json:"username"`
}

func CreateCustomer(db *pg.DB, username string) (*Users, error) {
	req := Users{
		Username: username,
	}

	_, err := db.Model(&req).Insert()
	if err != nil {
		return nil, err
	}

	user := &Users{}
	err = db.Model(user).
		Where("users.username = ?", username).
		Select()

	return user, err
}

func GetCustomer(db *pg.DB, username string) (*Users, error) {
	user := &Users{}
	err := db.Model(user).
		Where("users.username = ?", username).
		Select()

	return user, err
}

func GetCustomers(db *pg.DB) ([]*Users, error) {
	slugs := make([]*Users, 0)
	err := db.Model(&slugs).
		Select()

	return slugs, err
}

func UpdateCustomer(db *pg.DB, u *Users) (*Users, error) {
	_, err := db.Model(u).WherePK().Update()
	if err != nil {
		return nil, err
	}

	//user := &Users{}
	//err = db.Model(user).
	//	Relation("Users").
	//	Where("user.id = ?", u.ID).
	//	Select()

	return nil, err
}

func DeleteCustomer(db *pg.DB, userID int64) error {
	//user := &Users{
	//	ID: userID,
	//}
	//
	//err := db.Model(user).
	//	Relation("Users").
	//	Where("user.id = ?", user.ID).
	//	Select()
	//if err != nil {
	//	return err
	//}
	//
	//_, err = db.Model(user).WherePK().Delete()

	return nil
}
