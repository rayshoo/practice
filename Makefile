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