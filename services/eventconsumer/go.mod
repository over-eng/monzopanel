module github.com/over-eng/monzopanel/services/eventconsumer

go 1.22.7

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.6.0
	github.com/gocql/gocql v1.7.0
	github.com/over-eng/monzopanel/libraries/cassandratools v0.0.0-20241114115038-8d5a1ca9b284
	github.com/over-eng/monzopanel/libraries/kafkatools v0.0.0-20241114115038-8d5a1ca9b284
	github.com/over-eng/monzopanel/libraries/models v0.0.0-20241114115038-8d5a1ca9b284
	github.com/rs/zerolog v1.33.0
	go.mau.fi/zeroconfig v0.1.3
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.27.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)