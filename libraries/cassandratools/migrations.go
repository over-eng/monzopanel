package cassandratools

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
)

type Migration struct {
	Up      func(*gocql.Session) error
	Version string
	Down    func(*gocql.Session) error
}

func RunMigrations(ctx context.Context, session *gocql.Session, migrations []Migration) error {
	log := log.Ctx(ctx)
	err := assertMigrationsInOrder(migrations)
	if err != nil {
		return errors.New("migrations are not in order")
	}

	currentVersion, err := getCurrentVersion(session)
	if err != nil {
		return errors.New("could not determine the current migration version")
	}

	if currentVersion == getLatestVersion(migrations) {
		log.Info().Str("version", currentVersion).Msg("no migrations run, already on the latest version")
		return nil
	}

	for _, migration := range migrations {
		if migration.Version > currentVersion {
			log.Info().Str("version", migration.Version).Msg("running migration")
			err = runMigration(session, migration)
			if err != nil {
				return errors.New("failed to apply migration")
			}
		}
	}

	return nil
}

func runMigration(session *gocql.Session, migration Migration) error {
	err := migration.Up(session)
	if err != nil {
		return err
	}

	cql := "INSERT INTO migrations (version, applied_at) VALUES (?, ?)"
	err = session.Query(cql, migration.Version, time.Now()).Exec()
	if err != nil {
		return errors.Join(errors.New("failed to store current migration version in database"), err)
	}
	return nil
}

func assertMigrationsInOrder(migrations []Migration) error {
	for i, migration := range migrations {
		if i == 0 {
			continue
		}

		if migration.Version <= migrations[i-1].Version {
			return errors.New("migration invalid, versions must be alphabetical")
		}
	}
	return nil
}

func getLatestVersion(migrations []Migration) string {
	return migrations[len(migrations)-1].Version
}

func getCurrentVersion(session *gocql.Session) (string, error) {
	exists, err := migrationTableExists(session)
	if err != nil {
		return "", err
	}
	if !exists {
		err = createMigrationTable(session)
		if err != nil {
			return "", err
		}
	}

	// Note: the default ordering ensures we have the latest version
	cql := `
		SELECT version
		FROM migrations
		LIMIT 1
	`
	var version string
	err = session.Query(cql).Scan(&version)
	if err != nil && !errors.Is(err, gocql.ErrNotFound) {
		return "", err
	}

	return version, nil
}

func migrationTableExists(session *gocql.Session) (bool, error) {
	cql := `
        SELECT COUNT(*) 
        FROM system_schema.tables 
        WHERE keyspace_name = 'events' AND table_name = 'migrations'
    `

	var count int
	err := session.Query(cql).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createMigrationTable(session *gocql.Session) error {
	cql := `
		CREATE TABLE IF NOT EXISTS migrations (
			version text,
			applied_at timestamp,
			PRIMARY KEY (applied_at, version)
		) WITH CLUSTERING ORDER BY (version DESC)
    `

	err := session.Query(cql).Exec()
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	return nil
}
