BIN=droguedrums
VERSION=0.9
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"
SRCDIR=src
SRCFILES=$(wildcard ${SRCDIR}/*)

${BIN}: ${SRCFILES} Makefile
	go build ${LDFLAGS} -o ${BIN} ${SRCFILES}

test:
	cd ${SRCDIR} && go test

clean:
	rm ${BIN}
