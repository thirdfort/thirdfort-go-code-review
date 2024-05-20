.PHONY: build

BINARY     = build/consumer-api
VERSION    = 1.0.0
BUILD_TIME = `date -u '+%FT%TZ'`
GIT_HASH   = `git rev-parse HEAD`

SOURCEDIR        = .
BUILD_SOURCES   := $(shell find $(SOURCEDIR) -name '*.go')
LINT_SOURCES    := $(shell find $(SOURCEDIR) -name '*.go' | grep -v "/vendor/")
LDFLAGS          = -ldflags "-X github.com/thirdfort/thirdfort-go-code-review/internal.Version=$(VERSION) -X github.com/thirdfort/thirdfort-go-code-review/internal.BuildTime=$(BUILD_TIME) -X github.com/thirdfort/thirdfort-go-code-review/internal.GitHash=$(GIT_HASH)"
GOPATH           = $(shell go env GOPATH)

all: $(BINARY)

$(BINARY): $(BUILD_SOURCES)
	go build ${LDFLAGS} -o $(BINARY) main.go

run:
	# trap '$(MAKE) revoke' EXIT; go run main.go
	go run ${LDFLAGS} main.go

build:
	go build ${LDFLAGS} -o $(BINARY) main.go

vendor:
	go mod vendor
	cd testing/mock_services/platformapi && go mod vendor && cd ../../..

docker\:build: docker-up

docker-build: vendor
	CGO_ENABLED=0 go build ${LDFLAGS} -o consumer-api main.go

docker-up: vendor
	docker-compose up -d

docker-up-test: vendor
	docker-compose -f docker-compose-testing.yml up -d

docker\:stop: docker-stop

docker-stop:
	docker-compose down

docker\:clean: docker-clean

docker-clean:
	docker-compose down --rmi all

docker-clean-test:
	docker-compose -f docker-compose-testing.yml down --rmi all

docker\:logs:
	docker-compose logs -f

docker\:rebuild: docker-rebuild

docker-rebuild: docker-clean docker-up

docker-rebuild-test: docker-clean-test docker-up-test

docker\:run:
	docker run -it $(docker build -q .)

apidef: build
	@go run ./tools/gen_api_defs.go $(filter-out $@,$(MAKECMDGOALS))
	@echo "\n"

clean:
	go clean
	if [ -f $(BINARY) ] ; then rm $(BINARY); fi
	rm -rf ./vendor
	rm -rf ./testing/mock_services/platformapi/vendor
	for i in $(shell find . -name 'mock_*.go'); do rm $$i; done

generate:
	if [ ! -f ${GOPATH}/bin/mockgen ] ; then go install go.uber.org/mock/mockgen@latest; fi
	go generate gen.go

lint:
	gofmt -s -d -e ${LINT_SOURCES}

lint\:ci:
	test -z "$(shell gofmt -s -l ${LINT_SOURCES})"

lint\:fix:
	gofmt -s -l -w ${LINT_SOURCES}

test:
	go test -race -count=1 ${LDFLAGS} ./...

test\:integration: docker-rebuild-test
	until [ "`docker inspect -f {{.State.Status}} consumer-api-test`"=="running" ]; do \
		echo `docker inspect -f {{.State.Status}} consumer-api-test`; \
		sleep 1; \
	done;

	ENV=testing go test -tags=integration ./...

test\:ci: generate
	go test -coverprofile=coverage.out -covermode=atomic -race ${LDFLAGS} ./...

%:;
