package testtools

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
)

type CassandraSuite struct {
	container *cassandra.CassandraContainer
	Host      string
}

func NewCassandraSuite(ctx context.Context) (*CassandraSuite, error) {
	container, err := cassandra.Run(ctx, "cassandra:4.1.3")
	if err != nil {
		return nil, err
	}

	host, err := container.ConnectionHost(ctx)
	if err != nil {
		return nil, err
	}

	suite := &CassandraSuite{
		Host:      host,
		container: container,
	}

	return suite, nil
}

func (c *CassandraSuite) TearDownSuite() error {
	return testcontainers.TerminateContainer(c.container)
}
