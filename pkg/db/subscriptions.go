package db

import (
	"github.com/go-pg/pg/v10"
)

type Subscription struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	SlugID      string `json:"slug_id"`
	Condition   string `json:"condition"`
	TargetPrice int64  `json:"target_price"`
}

func CreateSubscriptions(db *pg.DB, req *Subscription) (*Subscription, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return nil, err
	}

	subscription := &Subscription{}
	err = db.Model(subscription).
		Relation("User").
		Where("subscription.id = ?", req.ID).
		Select()

	return subscription, err
}

func GetSubscription(db *pg.DB, subID string) (*Subscription, error) {
	subscription := &Subscription{}
	err := db.Model(subscription).
		Relation("User").
		Where("subscription.id = ?", subID).
		Select()

	return subscription, err
}

func GetSubscriptions(db *pg.DB) ([]*Subscription, error) {
	subscriptions := make([]*Subscription, 0)
	err := db.Model(&subscriptions).
		Relation("User").
		Select()

	return subscriptions, err
}

func UpdateSubscription(db *pg.DB, s *Subscription) (*Subscription, error) {
	_, err := db.Model(s).WherePK().Update()
	if err != nil {
		return nil, err
	}

	subscription := &Subscription{}
	err = db.Model(subscription).
		Relation("User").
		Where("subscription.id = ?", s.ID).
		Select()

	return subscription, err
}

func DeleteSubcription(db *pg.DB, slugID int64) error {
	subscription := &Subscription{
		ID: slugID,
	}

	err := db.Model(subscription).
		Relation("User").
		Where("subscription.id = ?", subscription.ID).
		Select()
	if err != nil {
		return err
	}

	_, err = db.Model(subscription).WherePK().Delete()

	return err
}
