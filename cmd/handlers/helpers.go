package handlers

import (
	"NFTracker/pkg/db"
	"NFTracker/pkg/opensea"
	"fmt"
	"github.com/go-pg/pg/v10"
	"log"
)

func closestMatchHelper(pgdb *pg.DB, slugQuery string) (string, error) {
	// Obtain namelist
	// TODO: optimise
	var namelist []string
	sluglist, err := db.GetAllSlugs(pgdb)
	if err != nil {
		log.Printf("[db.GetAllSlugs] Err: %v", err)
		return "", err
	}
	for _, slug := range sluglist {
		namelist = append(namelist, slug.SlugName)
	}

	// Find the closest match from list of popular NFTs from database
	matches := opensea.FindClosestmatch(slugQuery, 1, namelist)
	if matches[0] == "" {
		return "boredapeyachtclub", nil
	}

	return matches[0], nil
}

func popularCollectionHelper(pgdb *pg.DB, osResponse *opensea.OSResponse, slugQuery string) {
	// Update DB if collection is not already in database and is popular
	retSlug, err := db.GetSlug(pgdb, slugQuery)
	if err == pg.ErrNoRows && osResponse.Collection.Stats.NumOwners > PopularCollectionNumOwners {
		log.Printf(fmt.Sprintf("[db.GetSlug] New popular slug... adding slug to DB... %s", slugQuery))
		_, _ = db.CreateSlug(pgdb, &db.Slugs{
			SlugName:   slugQuery,
			FloorPrice: osResponse.Collection.Stats.FloorPrice,
		})
	}

	// If collection metadata is outdated, update it
	if retSlug.FloorPrice != osResponse.Collection.Stats.FloorPrice {
		updatedSlug, err := db.UpdateSlug(pgdb, &db.Slugs{
			SlugName:   osResponse.Collection.Slug,
			FloorPrice: osResponse.Collection.Stats.FloorPrice,
		})
		if err != nil {
			log.Printf(fmt.Sprintf("[db.UpdateSlug] Err: %s", err))
		}
		log.Printf(fmt.Sprintf("[db.UpdateSlug] Successfully updated slug: %s", updatedSlug))
	}
}
