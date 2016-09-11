BIN=		droguedrums
GITDIR=		${BIN}
VERSION=	1.1
BUILD_TIME=	$(shell date +%FT%T%z)
LDFLAGS=	-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"
MIDILIB=	portmidi
PREFIX=		/usr/local
OS=         $(shell uname)

${BIN}: *.go Makefile
	go build -tags ${MIDILIB} ${LDFLAGS} -o ${BIN}

test:
	go test -tags portmidi
ifeq ($(OS),Darwin)
	go test -tags coremidi
endif

clean:
	rm -f ${BIN}

checkfmt:
	@gofmt -d *.go

install: ${BIN}
	install ${BIN} ${PREFIX}/bin

www: mkdocs.yml
	mkdocs build --clean

zip:
	cd .. && zip -r ${GITDIR}-${VERSION}.zip ${GITDIR} --exclude "${GITDIR}/${BIN}" "${GITDIR}/site/*" "${GITDIR}/cinder/*" "*/.*"

binzip: clean ${BIN}
	cd .. && zip ${GITDIR}-${VERSION}-$(shell uname -s).zip ${GITDIR}/${BIN}

.PHONY: www
