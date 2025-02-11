brokers := 10.254.1.51:29291,10.254.1.51:29292,10.254.1.51:29293
topic := testtest
group := 0

up:
	docker-compose -f kafka/docker-compose.yaml up -d
	docker-compose -f rabbitmq/docker-compose.yaml up -d
.PHONY: up

kafka:
	docker-compose -f kafka/docker-compose.yaml up -d
.PHONY: kafka

rabbitmq:
	docker-compose -f rabbitmq/docker-compose.yaml up -d
.PHONY: rabbitmq

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
.PHONY: cp