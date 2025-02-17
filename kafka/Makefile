brokers := 10.254.1.51:29291,10.254.1.51:29292,10.254.1.51:29293
topic := test3
group := 0

# zookeeper or kraft
kafka_mode := kraft
#kafka_mode := zookeeper

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
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic}" \
	main.go
.PHONY: sp
# sarama group consumes
scg:
	cd sarama/consumer/group && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic} -X main.group=${group}" \
	main.go
.PHONY: scg
# sarama non group consumes
sc:
	cd sarama/consumer/nonGroup && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic} -X main.group=${group}" \
	main.go
.PHONY: sc