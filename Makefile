# Borrowed from:
# https://gist.github.com/turtlemonvh/38bd3d73e61769767c35931d8c70ccb4

BINDIR = bin
BINARY = tonight
VET_REPORT = vet.report
TEST_REPORT = tests-tonight.xml
GOARCH = amd64

# Setup the -ldflags option for go build here, interpolate the variable values
# LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
build-linux: clean test linux build-front setup
build-darwin: clean test darwin build-front setup

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -rf ${BINDIR}

test:
	cd tonight; \
	if ! hash go2xunit 2>/dev/null; then go get  github.com/tebeka/go2xunit; fi
	go test -v ./... 2>&1 | go2xunit -output ${TEST_REPORT}; \
	cd .. >/dev/null; \

	cd app; \
	npm run unit; \
	cd .. >/dev/null;

linux:
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINDIR}/${BINARY} ./tonight/cmd/main.go

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BINDIR}/${BINARY} ./tonight/cmd/main.go

build-front:
	cd app; \
	npm run build; \
	cd .. >/dev/null

setup:
	mkdir ${BINDIR}/bleve; \
	cp tonight/bleve/mapping.json ${BINDIR}/bleve/mapping.json

# windows:
# 	cd ${BUILD_DIR}; \
# 	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe . ; \
# 	cd - >/dev/null

.PHONY: linux darwin test fmt clean windows build-front build-linux build-darwin setup
