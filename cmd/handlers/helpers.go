package handlers

import (
	"NFTracker/pkg/db"
	"NFTracker/pkg/opensea"
	"github.com/go-pg/pg/v10"
	"log"
)

func closestMatchHelper(pgdb *pg.DB, slugQuery string) ([]string, error) {
	// Obtain namelist
	// TODO: optimise
	var namelist []string
	sluglist, err := db.GetAllSlugs(pgdb)
	if err != nil {
		log.Printf("[db.GetAllSlugs] Err: %v", err)
		return nil, err
	}
	for _, slug := range sluglist {
		namelist = append(namelist, slug.SlugName)
	}

	// Find the closest match from list of popular NFTs from database
	matches := opensea.FindClosestmatch(slugQuery, 1, namelist)
	return matches, nil
}
