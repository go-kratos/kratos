.PHONY: confluent/kafka/* confluent/zookeeper/* confluent/registry/* confluent/start confluent/stop fmt vet errcheck test test/create_kafka_topics dependencies dependencies/*


default: fmt vet errcheck test


# Confluent platform tasks

confluent/start: confluent/rest/start

confluent/stop: confluent/rest/stop confluent/registry/stop confluent/kafka/stop confluent/zookeeper/stop

# Download & extract tasks

confluent/confluent.tgz:
	mkdir -p confluent && wget http://packages.confluent.io/archive/1.0/confluent-1.0-2.10.4.tar.gz -O confluent/confluent.tgz

confluent/EXTRACTED: confluent/confluent.tgz
	tar xzf confluent/confluent.tgz -C confluent --strip-components 1 && mkdir confluent/logs && touch confluent/EXTRACTED

# Zookeeper tasks

confluent/zookeeper/start: confluent/EXTRACTED
	nohup confluent/bin/zookeeper-server-start confluent/etc/kafka/zookeeper.properties 2> confluent/logs/zookeeper.err > confluent/logs/zookeeper.out < /dev/null &
	while ! nc localhost 2181 </dev/null; do echo "Waiting for zookeeper..."; sleep 1; done

confluent/zookeeper/stop: confluent/EXTRACTED
	confluent/bin/zookeeper-server-stop

# Kafka tasks

confluent/kafka/start: confluent/zookeeper/start confluent/EXTRACTED
	nohup confluent/bin/kafka-server-start confluent/etc/kafka/server.properties 2> confluent/logs/kafka.err > confluent/logs/kafka.out < /dev/null &
	while ! nc localhost 9092 </dev/null; do echo "Waiting for Kafka..."; sleep 1; done

confluent/kafka/stop: confluent/EXTRACTED
	confluent/bin/kafka-server-stop

# schema-registry tasks

confluent/registry/start: confluent/kafka/start confluent/EXTRACTED
	nohup confluent/bin/schema-registry-start confluent/etc/schema-registry/schema-registry.properties 2> confluent/logs/schema-registry.err > confluent/logs/schema-registry.out < /dev/null &
	while ! nc localhost 8081 </dev/null; do echo "Waiting for schema registry..."; sleep 1; done

confluent/registry/stop: confluent/EXTRACTED
	confluent/bin/kafka-server-stop

# REST proxy tasks

confluent/rest/start: confluent/registry/start confluent/EXTRACTED
	nohup confluent/bin/kafka-rest-start confluent/etc/kafka-rest/kafka-rest.properties 2> confluent/logs/kafka-rest.err > confluent/logs/kafka-rest.out < /dev/null &
	while ! nc localhost 8082 </dev/null; do echo "Waiting for REST proxy..."; sleep 1; done

confluent/rest/stop: confluent/EXTRACTED
	confluent/bin/kafka-rest-stop


# CI tasks

test:
	go test -v -race ./...

vet:
	go vet ./...

errcheck:
	errcheck ./...

fmt:
	@if [ -n "$$(go fmt ./...)" ]; then echo 'Please run go fmt on your code.' && exit 1; fi

dependencies: dependencies/errcheck dependencies/get

dependencies/errcheck:
	go get github.com/kisielk/errcheck

dependencies/get:
	go get -t ./...

test/create_kafka_topics: confluent/kafka/start
	confluent/bin/kafka-topics --create --partitions  1 --replication-factor 1 --topic test.1  --zookeeper localhost:2181
	confluent/bin/kafka-topics --create --partitions  4 --replication-factor 1 --topic test.4  --zookeeper localhost:2181 --config retention.ms=604800000
	confluent/bin/kafka-topics --create --partitions 64 --replication-factor 1 --topic test.64 --zookeeper localhost:2181
