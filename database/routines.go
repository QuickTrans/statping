package database

import (
	"fmt"
	"github.com/hunterlong/statping/types"
	"github.com/hunterlong/statping/utils"
	"os"
	"time"
)

var (
	log               = utils.Log
	removeRowsAfter   = types.Month * 6
	maintenceDuration = types.Hour
)

func StartMaintenceRoutine() {
	dur := os.Getenv("REMOVE_AFTER")
	var removeDur time.Duration

	if dur != "" {
		parsedDur, err := time.ParseDuration(dur)
		if err != nil {
			log.Errorf("could not parse duration: %s, using default: %s", dur, removeRowsAfter.String())
			removeDur = removeRowsAfter
		} else {
			removeDur = parsedDur
		}
	} else {
		removeDur = removeRowsAfter
	}

	log.Infof("Service Failure and Hit records will be automatically removed after %s", removeDur.String())
	go databaseMaintence(removeDur)
}

// databaseMaintence will automatically delete old records from 'failures' and 'hits'
// this function is currently set to delete records 7+ days old every 60 minutes
func databaseMaintence(dur time.Duration) {
	deleteAfter := time.Now().UTC().Add(dur)

	for range time.Tick(maintenceDuration) {
		log.Infof("Deleting failures older than %s", dur.String())
		DeleteAllSince("failures", deleteAfter)

		log.Infof("Deleting hits older than %s", dur.String())
		DeleteAllSince("hits", deleteAfter)

		maintenceDuration = types.Hour
	}
}

// DeleteAllSince will delete a specific table's records based on a time.
func DeleteAllSince(table string, date time.Time) {
	sql := fmt.Sprintf("DELETE FROM %v WHERE created_at < '%v';", table, database.FormatTime(date))
	db := database.Exec(sql)
	if db.Error() != nil {
		log.Warnln(db.Error())
	}
}