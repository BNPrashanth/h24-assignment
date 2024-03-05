VERSION = 0.0.1
CONTAINER = h24-analyser-service

PORT = 8080
BASE_URL = http://localhost:${PORT}

HAS_EXITED := $(shell docker ps -a | grep ${CONTAINER})

run-debug:
	@echo "+ $@"
	@export BASE_URL=${BASE_URL}
	@export PORT=${PORT}
	@go run cmd/main.go

build:
	@echo "+ $@"
	@docker build --pull -t ${CONTAINER}:$(VERSION) .

rm:
ifdef HAS_EXITED
	@echo "+ $@"
	@docker rm -f ${CONTAINER}
endif

run: rm build
	@echo "+ $@"
	@docker run --name ${CONTAINER} \
		-p ${PORT}:${PORT} \
		-e "PORT=$(PORT)" \
		-e "BASE_URL=$(BASE_URL)" \
		-d ${CONTAINER}:$(VERSION)
	@sleep 1
	@docker logs ${CONTAINER}
