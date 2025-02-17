fmt:
	goimports -l -w .
.PHONY: fmt

avro:
	cd avro && \
	go run main.go
.PHONY: avro

k-up:
	$(MAKE) -C kafka up
.PHONY: k-up
k-down:
	$(MAKE) -C kafka down
.PHONY: k-down

r-up:
	$(MAKE) -C rabbitmq up
.PHONY: r-up
r-down:
	$(MAKE) -C rabbitmq down
.PHONY: r-down

wasm-up:
	$(MAKE) -C webassembly up
.PHONY: wasm-up

wasm-build:
	$(MAKE) -C webassembly build
.PHONY: wasm-build