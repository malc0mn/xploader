SOURCEDIR=cmd
BINARY=xpinfo

LDFLAGS=-ldflags "-s -w"

.DEFAULT_GOAL: all

.PHONY: all
all: xpinfo

xpinfo:
	cd ${SOURCEDIR}; go build -trimpath ${LDFLAGS} -o ../${BINARY}

.PHONY: test
test:
	go test ./...

.PHONY: testv
testv:
	go test -v ./...

.PHONY: install
install:
	cd ${SOURCEDIR}; GOBIN=/usr/local/bin/ go install ${LDFLAGS}

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
