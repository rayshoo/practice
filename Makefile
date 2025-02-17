brokers := 10.254.1.51:29291,10.254.1.51:29292,10.254.1.51:29293
topic := test3
group := 0

up:
	docker-compose -f kafka/docker-compose.yaml up -d
	docker-compose -f rabbitmq/docker-compose.yaml up -d
.PHONY: up
down:
	docker-compose -f kafka/docker-compose.yaml down
	docker-compose -f rabbitmq/docker-compose.yaml down
.PHONY: down

k-up:
	docker-compose -f kafka/docker-compose.yaml up -d
.PHONY: k-up
k-down:
	docker-compose -f kafka/docker-compose.yaml down -d
.PHONY: k-down

r-up:
	docker-compose -f rabbitmq/docker-compose.yaml up -d
.PHONY: r-up
r-down:
	docker-compose -f rabbitmq/docker-compose.yaml down -d
.PHONY: r-down

fmt:
	goimports -l -w .
.PHONY: fmt

# confluent produce
cp:
	cd kafka/confluent/producer && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic}" \
	main.go
.PHONY: cp
# confluent consumes
cc:
	cd kafka/confluent/consumer && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic} -X main.group=${group}" \
	main.go
.PHONY: cc

# sarama produce
sp:
	cd kafka/sarama/producer && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic}" \
	main.go
.PHONY: sp
# sarama group consumes
scg:
	cd kafka/sarama/consumer/group && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic} -X main.group=${group}" \
	main.go
.PHONY: scg
# sarama non group consumes
sc:
	cd kafka/sarama/consumer/nonGroup && \
	go run \
	-ldflags "-s -w -X main.brokers=${brokers} -X main.topic=${topic} -X main.group=${group}" \
	main.go
.PHONY: sc

avro:
	cd avro && \
	go run main.go
.PHONY: avro