package eventstore

import (
	"github.com/gocql/gocql"
	"github.com/over-eng/monzopanel/libraries/cassandratools"
)

var migrations = []cassandratools.Migration{
	_01_create_events,
}

var _01_create_events = cassandratools.Migration{
	Version: "0001_create_events",
	Up: func(session *gocql.Session) error {
		cql := `
			CREATE TABLE IF NOT EXISTS events (
				id text,
				event text,
				team_id text,
				distinct_id text,
				properties text,
				client_timestamp timestamp,
				created_at timestamp,
				loaded_at timestamp,
				PRIMARY KEY (created_at, team_id, event, distinct_id, id)
			)
		`
		return session.Query(cql).Exec()
	},
	Down: func(session *gocql.Session) error {
		return session.Query("DROP TABLE IF EXISTS events").Exec()
	},
}
