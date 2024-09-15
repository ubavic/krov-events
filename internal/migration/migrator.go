package migration

import (
	"database/sql"
	"embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

type migration struct {
	id   int
	code string
}

type Migrator struct {
	db         *sql.DB
	migrations []migration
}

func NewMigrator(db *sql.DB, migrationFS *embed.FS) (*Migrator, error) {
	migrations := []migration{}

	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return nil, nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			segments := strings.Split(entry.Name(), ".")
			if len(segments) != 2 {
				continue
			}

			id, err := strconv.Atoi(segments[0])
			if err != nil {
				continue
			}

			fileContent, err := migrationFS.ReadFile("migrations/" + entry.Name())
			if err != nil {
				return nil, err
			}

			migrations = append(migrations, migration{
				code: string(fileContent),
				id:   id,
			})
		}
	}

	sortFunc := func(a migration, b migration) int {
		return a.id - b.id
	}

	slices.SortFunc(migrations, sortFunc)

	migrator := Migrator{
		db:         db,
		migrations: migrations,
	}

	return &migrator, nil
}

func (migrator *Migrator) RunMigrations() error {
	for _, migration := range migrator.migrations {
		start := time.Now()
		_, err := migrator.db.Exec(migration.code)
		end := time.Now()

		if err != nil {
			fmt.Println(err, start, end)
		}
	}

	return nil
}
