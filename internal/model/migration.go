package model

import "time"

type MigrationId string

type Migration struct {
	Id      MigrationId
	Time    time.Time
	Success bool
}
