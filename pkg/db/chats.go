package db

import (
	"github.com/go-pg/pg/v10"
)

type Chats struct {
	Id          int64  `pg:",pk" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Invitelink  string `json:"invitelink"`
}

func CreateChat(db *pg.DB, req *Chats) (*Chats, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return nil, err
	}

	chat := &Chats{}
	err = db.Model(chat).
		Where("chats.id = ?", req.Id).
		Select()

	return chat, err
}

func GetChat(db *pg.DB, id int64) (*Chats, error) {
	chat := &Chats{}
	err := db.Model(chat).
		Where("chats.id = ?", id).
		Select()

	return chat, err
}

func UpdateChat(db *pg.DB, s *Chats) (*Chats, error) {
	_, err := db.Model(s).WherePK().Update()
	if err != nil {
		return nil, err
	}

	chat := &Chats{}
	err = db.Model(chat).
		Where("chats.id = ?", s.Id).
		Select()

	return chat, err
}
