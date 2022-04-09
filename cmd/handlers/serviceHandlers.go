package handlers

import (
	"fmt"
	"strconv"
	"time"
)

// Returns a channel, to which a carpark's details are
// written each time data is received in which more
// than the specified number of places are available
func alertSpaces(a chan<- string, c int) (e error) {
	// Whilst this is implemented using ksqlDB, perhaps it could be done
	// using a normal Kafka consumer and use Golang to apply the filter
	//
	// TODO01 add a channel so that user can run another
	// command to delete an alert
	//
	// TODO02 return a chan of carPark type, not a hardcoded string
	//
	defer close(a)

	// Run this for five minutes and then exit
	// TODO: does this actually do this?
	const queryResultTimeoutSeconds = 300

	// Prepare the request
	k := "SELECT NAME, TS, EMPTY_PLACES, CAPACITY"
	k += "  FROM CARPARK_EVENTS"
	k += " WHERE  EMPTY_PLACES > " + strconv.Itoa(c)
	k += " EMIT CHANGES;"

	// This Go routine will handle rows as and when they
	// are sent to the channel
	rc := make(chan ksqldb.Row)
	hc := make(chan ksqldb.Header, 1)

	var CARPARK string
	var DATA_TS float64
	var CURRENT_EMPTY_PLACES float64
	var CAPACITY float64
	go func() {

		for row := range rc {
			if row != nil {

				// Store the row values in the carPark object

				CARPARK = row[0].(string)
				DATA_TS = row[1].(float64)
				CURRENT_EMPTY_PLACES = row[2].(float64)
				CAPACITY = row[3].(float64)
				// Handle the timestamp
				t := int64(DATA_TS)
				ts := time.Unix(t/1000, 0).Format(time.RFC822)
				a <- fmt.Sprintf("✨ 🎉  🚗 The %v carpark has %v spaces available (capacity %v)\n(data as of %v)", CARPARK, CURRENT_EMPTY_PLACES, CAPACITY, ts)
			}
		}
	}()

	// Do the request
	ctx, cancel := context.WithTimeout(context.TODO(), queryResultTimeoutSeconds*time.Second)
	defer cancel()

	client := ksqldb.NewClient(KSQLDB_ENDPOINT, KSQLDB_API_KEY, KSQLDB_API_SECRET).Debug()

	e = client.Push(ctx, k, "latest", rc, hc)

	if e != nil {
		return fmt.Errorf("Error running Push request against ksqlDB:\n%v", e)
	}

	return nil
}
