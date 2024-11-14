package cassandratools_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/over-eng/monzopanel/libraries/cassandratools"
	"github.com/over-eng/monzopanel/libraries/testtools"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	cassandra *testtools.CassandraSuite
}

func (suite *testSuite) SetupSuite() {
	cassandrasuite, err := testtools.NewCassandraSuite(context.Background())
	suite.Require().NoError(err)
	suite.cassandra = cassandrasuite
}

func (suite *testSuite) TearDownSuite() {
	err := suite.cassandra.TearDownSuite()
	suite.Require().NoError(err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) TestStartSession() {

	cfg := cassandratools.ConnectionConfig{
		Hosts:       []string{suite.cassandra.Host},
		User:        "cassandra",
		Password:    "cassandra",
		Consistency: "one",
	}

	keyspace := cassandratools.Keyspace{
		Name:              "test_keyspace",
		Class:             "SimpleStrategy",
		ReplicationFactor: 1,
	}

	migrations := []cassandratools.Migration{
		{
			Version: "001_create_test_table",
			Up: func(session *gocql.Session) error {
				cql := `
					CREATE TABLE test_table (
						id text,
						description text,
						PRIMARY KEY (id)
					)
				`
				return session.Query(cql).Exec()
			},
			Down: func(session *gocql.Session) error {
				return session.Query("DROP TABLE test_table").Exec()
			},
		},
		{
			Version: "002_create_test_table2",
			Up: func(session *gocql.Session) error {
				cql := `
					CREATE TABLE test_table2 (
						id text,
						description text,
						PRIMARY KEY (id)
					)
				`
				return session.Query(cql).Exec()
			},
			Down: func(session *gocql.Session) error {
				return session.Query("DROP TABLE test_table2").Exec()
			},
		},
	}

	suite.Run("create a simple session", func() {
		session, err := cassandratools.NewSession(cfg).Start(context.Background())
		suite.Assert().NoError(err)
		defer session.Close()

		err = session.Query("SELECT keyspace_name FROM system_schema.keyspaces").Exec()
		suite.Assert().NoError(err, "should be able to query with session")
	})

	suite.Run("with keyspace creation", func() {
		session, err := cassandratools.NewSession(cfg).
			WithCreateKeyspace(keyspace).
			Start(context.Background())
		suite.Assert().NoError(err)

		replicationMap := make(map[string]string)
		err = session.
			Query("SELECT replication FROM system_schema.keyspaces WHERE keyspace_name = ?", keyspace.Name).
			Scan(&replicationMap)
		suite.Assert().NoError(err)
		suite.Assert().Equal("org.apache.cassandra.locator.SimpleStrategy", replicationMap["class"])
		suite.Assert().Equal("1", replicationMap["replication_factor"])

		session.Close()

		session, err = cassandratools.NewSession(cfg).
			WithCreateKeyspace(keyspace).
			Start(context.Background())
		suite.Assert().NoError(err, "calling create keyspace a second time shouldn't cause an error")
		defer session.Close()
		err = session.Query(fmt.Sprintf("DROP keyspace %s", keyspace.Name)).Exec()
		suite.Assert().NoError(err)
	})

	suite.Run("migrations are applied in order", func() {
		session, err := cassandratools.NewSession(cfg).
			WithCreateKeyspace(keyspace).
			WithMigrations(migrations).
			WithUseKeyspace(keyspace.Name).
			Start(context.Background())
		suite.Assert().NoError(err)
		defer session.Close()

		// assert the order of the migrations table
		iter := session.Query("SELECT version, applied_at FROM migrations ORDER BY applied_at ASC").Iter()

		var version string
		versions := []string{}

		var appliedAt time.Time
		appliedAts := []time.Time{}
		for {
			if !iter.Scan(&version, &appliedAt) {
				break
			}
			versions = append(versions, version)
			appliedAts = append(appliedAts, appliedAt)
		}
		for i := 0; i < len(versions)-1; i++ {
			suite.Assert().True(versions[i+1] > versions[i])
			suite.Assert().True(appliedAts[i+1].UnixMilli() > appliedAts[i].UnixMilli())
		}
		err = session.Query(fmt.Sprintf("DROP keyspace %s", keyspace.Name)).Exec()
		suite.Assert().NoError(err)
	})

	suite.Run("migrations only apply unseen versions", func() {
		// run just the first migration
		firstMigration := []cassandratools.Migration{migrations[0]}
		session, err := cassandratools.NewSession(cfg).
			WithCreateKeyspace(keyspace).
			WithMigrations(firstMigration).
			WithUseKeyspace(keyspace.Name).
			Start(context.Background())
		suite.Assert().NoError(err)
		session.Close()

		// now apply all of them, the first migration is written so that it should
		// fail if applied twice, so if we have no errors then we are good.
		session, err = cassandratools.NewSession(cfg).
			WithCreateKeyspace(keyspace).
			WithMigrations(migrations).
			WithUseKeyspace(keyspace.Name).
			Start(context.Background())
		suite.Assert().NoError(err)
		defer session.Close()

		err = session.Query(fmt.Sprintf("DROP keyspace %s", keyspace.Name)).Exec()
		suite.Assert().NoError(err)
	})

	suite.Run("out of order migrations fail", func() {
		outOfOrderMigrations := []cassandratools.Migration{migrations[1], migrations[0]}
		_, err := cassandratools.NewSession(cfg).
			WithCreateKeyspace(keyspace).
			WithMigrations(outOfOrderMigrations).
			WithUseKeyspace(keyspace.Name).
			Start(context.Background())
		suite.Assert().Error(err)
	})

	suite.Run("migrations error when keyspace isn't set", func() {
		_, err := cassandratools.NewSession(cfg).
			WithCreateKeyspace(keyspace).
			WithMigrations(migrations).
			Start(context.Background())
		suite.Assert().Error(err)
	})
}
