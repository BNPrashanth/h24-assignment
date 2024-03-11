VERSION = 0.0.1
CONTAINER = h24-analyser-service

PORT = 8080
BASE_URL = http://localhost:${PORT}

HAS_EXITED := $(shell docker ps -a | grep ${CONTAINER})

BUILD_EXISTS := $(shell ls -a | grep ${CONTAINER})

run-debug:
	@echo "+ $@"
	export BASE_URL=${BASE_URL} && \
	export PORT=${PORT} && \
	go run cmd/main.go

build:
	@echo "+ $@"
	@docker build --pull -t ${CONTAINER}:$(VERSION) .

rm:
ifdef HAS_EXITED
	@echo "+ $@"
	@docker rm -f ${CONTAINER}
endif

test:
	@echo "+ $@"
	go test -v ./...

run: rm test build
	@echo "+ $@"
	@docker run --name ${CONTAINER} \
		-p ${PORT}:${PORT} \
		-e "PORT=$(PORT)" \
		-e "BASE_URL=$(BASE_URL)" \
		-d ${CONTAINER}:$(VERSION)
	@sleep 1
	@docker logs ${CONTAINER}

go-clean:
ifdef BUILD_EXISTS
	@echo "+ $@"
	go clean
	rm ${CONTAINER}
endif

go-build:
	@echo "+ $@"
	GOARCH=amd64 GOOS=darwin go build -o ${CONTAINER} cmd/main.go
 	GOARCH=amd64 GOOS=linux go build -o ${CONTAINER} cmd/main.go
 	GOARCH=amd64 GOOS=windows go build -o ${CONTAINER} cmd/main.go

go-run: go-clean test go-build
	export BASE_URL=${BASE_URL} && \
	export PORT=${PORT} && \
	./${CONTAINER}
