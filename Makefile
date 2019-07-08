BINARY=$(shell pwd | sed -e "s/.*\///")
GOARCH=amd64
UNAME=$(shell uname -s)

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
BUILD_DIR=${GOPATH}/src/${BINARY}
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})
BIN_DIR=${BUILD_DIR}/bin
SOURCEDIR=${BUILD_DIR}/src
VET_REPORT = ${BUILD_DIR}/vet.report
TEST_REPORT = ${BUILD_DIR}/tests.xml

NODEMON_BIN=nodemon

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
all: deps clean test vet fmt linux darwin windows

build: deps clean vet
build:
ifeq ($(OS),Windows_NT)
	make windows
else ifeq ($(UNAME),Darwin)
	make darwin
else ifeq ($(UNAME),Linux)
	make linux
endif

run: deps
	cd ${BUILD_DIR}; \
	go run -race ${SOURCEDIR}

watch: deps
	cd ${BUILD_DIR}; \
	$(NODEMON_BIN) -e go -x 'go run -race ${SOURCEDIR} || exit 1'

deps: vendor
vendor: link
	cd ${BUILD_DIR}; \
	dep ensure

	touch $@

link:
	if [ "${BUILD_DIR_LINK}" != "${CURRENT_DIR}" ] && [ "${BUILD_DIR}" != "${CURRENT_DIR}" ]; then \
		echo "Fixing symlinks for build"; \
		rm -f ${BUILD_DIR}; \
		ln -s ${CURRENT_DIR} ${BUILD_DIR}; \
	fi

linux:
	cd ${SOURCEDIR}; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY}-linux-${GOARCH} .

darwin:
	cd ${SOURCEDIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY}-darwin-${GOARCH} .

windows:
	cd ${SOURCEDIR}; \
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY}-windows-${GOARCH}.exe .

# test:
# 	export $$(cat ./.env.test | xargs); \
# 	cd ${BUILD_DIR}; \
# 	cd $(shell find . -type d -name "go2xunit"); \
# 	go build;
#
# 	export $$(cat ./.env.test | xargs); \
# 	cd ${SOURCEDIR}; \
# 	go test -v ./... 2>&1 | ${BUILD_DIR}/$(shell find . -type f -name "go2xunit") -output ${TEST_REPORT} ; \
# 	cd - >/dev/null

test:
	export $$(cat ./.env.test | xargs); \
	cd ${SOURCEDIR}; \
	go test -race -v ./...

vet:
	cd ${SOURCEDIR}; \
	go vet ./... > ${VET_REPORT} 2>&1

fmt:
	cd ${SOURCEDIR}; \
	go fmt $$(go list ./... | grep -v /vendor/)

clean:
	go clean
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BIN_DIR}/${BINARY}-*

.PHONY: link linux darwin windows test vet fmt clean
