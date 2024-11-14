package cassandratools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gocql/gocql"
)

type sessionBuilder struct {
	cfg             ConnectionConfig
	migrations      []Migration
	createKeyspace  Keyspace
	appliedKeyspace string
}

type ConnectionConfig struct {
	Hosts       []string `yaml:"hosts"`
	User        string   `yaml:"user"`
	Password    string   `yaml:"password"`
	Consistency string   `yaml:"consistency"`
}

type Keyspace struct {
	Name              string `yaml:"name"`
	Class             string `yaml:"class"`
	ReplicationFactor int    `yaml:"replication_factor"`
}

func NewSession(cfg ConnectionConfig) *sessionBuilder {
	s := &sessionBuilder{cfg: cfg}
	return s
}

func (sb *sessionBuilder) WithCreateKeyspace(keyspace Keyspace) *sessionBuilder {
	sb.createKeyspace = keyspace
	return sb
}

func (sb *sessionBuilder) WithUseKeyspace(keyspaceName string) *sessionBuilder {
	sb.appliedKeyspace = keyspaceName
	return sb
}

func (sb *sessionBuilder) WithMigrations(migrations []Migration) *sessionBuilder {
	sb.migrations = migrations
	return sb
}

func (sb *sessionBuilder) Start(ctx context.Context) (*gocql.Session, error) {
	cluster := gocql.NewCluster(sb.cfg.Hosts...)
	cluster.Consistency = gocql.ParseConsistency(sb.cfg.Consistency)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: sb.cfg.User,
		Password: sb.cfg.Password,
	}

	if sb.createKeyspace.Name != "" {
		err := createKeyspaceIfNotExists(cluster, sb.createKeyspace)
		if err != nil {
			return nil, err
		}
	}

	if sb.appliedKeyspace != "" {
		cluster.Keyspace = sb.appliedKeyspace
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	if sb.migrations != nil {

		if sb.appliedKeyspace == "" {
			return nil, errors.New("migrations need to be scoped to a keyspace, use WithUseKeyspace to run them")
		}

		err = runMigrations(ctx, session, sb.migrations, sb.appliedKeyspace)
		if err != nil {
			session.Close()
			return nil, err
		}

	}
	return session, nil
}

func createKeyspaceIfNotExists(cluster *gocql.ClusterConfig, keyspace Keyspace) error {
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if keyspace.Name == "" {
		return errors.New("no keyspace name provided")
	}

	if keyspace.Class == "" {
		return errors.New("no replication class provided")
	}

	if keyspace.ReplicationFactor <= 0 {
		return errors.New("no replication factor provided")
	}

	cql := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s
		WITH replication = {'class': '%s', 'replication_factor' : %d}
	`, keyspace.Name, keyspace.Class, keyspace.ReplicationFactor)
	return session.Query(cql).Exec()
}
