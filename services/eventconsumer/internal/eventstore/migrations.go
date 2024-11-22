package eventstore

import (
	"github.com/gocql/gocql"
	"github.com/over-eng/monzopanel/libraries/cassandratools"
)

var migrations = []cassandratools.Migration{
	_01_create_events,
	_02_create_events_by_hour_counter,
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
				PRIMARY KEY (team_id, distinct_id, created_at, id)
			) WITH CLUSTERING ORDER BY (distinct_id ASC, created_at DESC)
		`
		return session.Query(cql).Exec()
	},
	Down: func(session *gocql.Session) error {
		return session.Query("DROP TABLE IF EXISTS events").Exec()
	},
}

var _02_create_events_by_hour_counter = cassandratools.Migration{
	Version: "0002_create_events_by_hour_counter",
	Up: func(session *gocql.Session) error {
		cql := `
			CREATE TABLE events_by_hour_counter (
				team_id text,
				distinct_id text,
				bucket_hour timestamp,
				event text,
				event_count counter,
				PRIMARY KEY ((team_id, bucket_hour), event, distinct_id)
			)
		`
		return session.Query(cql).Exec()
	},
	Down: func(session *gocql.Session) error {
		return session.Query("DROP TABLE IF EXISTS events_by_hour_counter").Exec()
	},
}
