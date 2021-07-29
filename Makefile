BINARY=engine
test: 
	go test -v -cover -covermode=atomic ./...

engine:
	go build -o ${BINARY} app/*.go


unittest:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker:
	docker build -t go-clean-arch .

run:
	docker-compose up --build -d

stop:
	docker-compose down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

.PHONY: clean install unittest build docker run stop vendor lint-prepare lint

# Definitions
ROOT                    := $(PWD)
GO_HTML_COV             := ./coverage.html
GO_TEST_OUTFILE         := ./c.out
GOLANG_DOCKER_IMAGE     := golang:1.15
CC_TEST_REPORTER_ID		:= ${CC_TEST_REPORTER_ID}
CC_PREFIX				:= github.com/khihadysucahyo/go-clean-arch-boilerplate

# custom logic for code climate, gross but necessary
coverage:
	# download CC test reported
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} \
		/bin/bash -c \
		"curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter"
	
	# update perms
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} chmod +x ./cc-test-reporter

	# run before build
	docker run -w /app -v ${ROOT}:/app \
		 -e CC_TEST_REPORTER_ID=${CC_TEST_REPORTER_ID} \
		${GOLANG_DOCKER_IMAGE} ./cc-test-reporter before-build

	# run testing
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} go test ./... -coverprofile=${GO_TEST_OUTFILE}
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} go tool cover -html=${GO_TEST_OUTFILE} -o ${GO_HTML_COV}

	#upload coverage result
	$(eval PREFIX=${CC_PREFIX})
ifdef prefix
	$(eval PREFIX=${prefix})
endif
	# upload data to CC
	docker run -w /app -v ${ROOT}:/app \
		-e CC_TEST_REPORTER_ID=${CC_TEST_REPORTER_ID} \
		${GOLANG_DOCKER_IMAGE} ./cc-test-reporter after-build --prefix ${PREFIX}

test:
	@go test ./... -coverprofile=./coverage.out & go tool cover -html=./coverage.out
