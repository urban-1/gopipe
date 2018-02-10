KAFKA_VERSION := v0.11.3
GO_FILES := $(shell find . -iname '*.go' -type f | grep -v /build/)
MODS :=./input ./proc ./output

.PHONY: all gopipe librdkafka fmt tests show_coverage rdkafka help


help:
	@echo ""
	@echo "Available targets:"
	@echo
	@echo "  rdkafka         Build librdkafka locally"
	@echo "  gopipe          Build gopipe"
	@echo "  tests           Run tests (includes fmt and vet)"
	@echo "  show_coverage   Show coverage results"
	@echo

all: rdkafka gopipe

gopipe:
	go build

rdkafka:
	-@mkdir -p build/src
	-@mkdir -p build/local
	@(  cd build/src; \
		git clone https://github.com/edenhill/librdkafka.git; \
		cd librdkafka; \
		git fetch origin && \
		git checkout $(KAFKA_VERSION) && \
		./configure --prefix=$(CURDIR)/build/local && \
		make && \
		make install \
	)

fmt:
	@go fmt ./...


vet:
	@go vet ./...

tests: fmt vet
	-@mkdir -p build/coverage
	@go get -u github.com/wadey/gocovmerge
	(   export PKG_CONFIG_PATH=$(CURDIR)/build/local/lib/pkgconfig; \
		export LD_LIBRARY_PATH=$(CURDIR)/build/local/lib; \
		LD_PRELOAD=$(CURDIR)/build/local/lib/librdkafka.so.1 go build github.com/confluentinc/confluent-kafka-go/kafka; \
		for mod in $(MODS); do \
			go test -v -cover -coverprofile=build/coverage/$$mod.out ./$$mod || exit 1; \
		done;\
	)
	@gocovmerge build/coverage/* > build/coverage/all.out

show_coverage:
	@go tool cover -html=build/coverage/all.out
