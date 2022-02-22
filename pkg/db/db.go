package db

import (
	"fmt"
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"log"
)

func NewDB() (*pg.DB, error) {
	// connect to db
	db := pg.Connect(&pg.Options{
		Addr:     "db:5432",
		User:     "postgres",
		Password: "admin",
	})

	// run migrations
	log.Printf("@@@@@ run migrations")
	collection := migrations.NewCollection()
	err := collection.DiscoverSQLMigrations("migrations")
	if err != nil {
		return nil, err
	}

	// init db
	log.Printf("@@@@@ init")
	_, _, err = collection.Run(db, "init")
	if err != nil {
		return nil, err
	}

	log.Printf("@@@@@ up")
	oldVersion, newVersion, err := migrations.Run(db, "up")
	if err != nil {
		return nil, err
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}

	//return the db connections
	return db, nil
}
