brokers := bootstrap.ezarrow.co.kr:443
topic := test
group := 0
cert := ..\..\files\ssl\server.crt

# zookeeper or kraft
#kafka_mode := kraft
kafka_mode := zookeeper

up:
	docker-compose -f deploy/docker/${kafka_mode}/docker-compose.yaml up -d
.PHONY: up
down:
	docker-compose -f deploy/docker/${kafka_mode}/docker-compose.yaml down
.PHONY: down

# confluent produce
cp:
	cd confluent/producer && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic}" \
	main.go
.PHONY: cp
# confluent consumes
cc:
	cd confluent/consumer && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic} -X main.group=${group}" \
	main.go
.PHONY: cc

# sarama produce
sp:
	cd sarama/producer && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic} -X main.cert=${cert}" \
	main.go
.PHONY: sp
# sarama group consumes
scg:
	cd sarama/consumer/group && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic} -X main.group=${group} -X main.cert=..\${cert}" \
	main.go
.PHONY: scg
# sarama non group consumes
sc:
	cd sarama/consumer/nonGroup && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic} -X main.group=${group} -X main.cert=..\${cert}" \
	main.go
.PHONY: sc