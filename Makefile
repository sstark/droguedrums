BIN=		droguedrums
VERSION=	0.9
BUILD_TIME=	$(shell date +%FT%T%z)
LDFLAGS=	-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"
SRC=		src
GOFILES:=	$(wildcard ${SRC}/*.go)
SRCFILES:=	$(patsubst ${SRC}/%_test.go,,${GOFILES})

${BIN}: ${SRCFILES} Makefile
	go build ${LDFLAGS} -o ${BIN} ${SRCFILES}

test:
	cd ${SRC} && go test

clean:
	rm ${BIN}
