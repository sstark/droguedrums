BIN=		droguedrums
GITDIR=		${BIN}
VERSION=	1.0
BUILD_TIME=	$(shell date +%FT%T%z)
LDFLAGS=	-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"
MIDILIB=	portmidi
PREFIX=		/usr/local

${BIN}: src/* Makefile
	@echo building for ${MIDILIB}
	cd src && go build -tags ${MIDILIB} ${LDFLAGS} -o ../${BIN}

test:
	cd src && go test -tags portmidi
	cd src && go test -tags coremidi

clean:
	rm -f ${BIN}

checkfmt:
	@gofmt -d src/*.go

install: ${BIN}
	install ${BIN} ${PREFIX}/bin

www: mkdocs.yml
	mkdocs build --clean

zip:
	cd .. && zip -r ${GITDIR}-${VERSION}.zip ${GITDIR} --exclude "${GITDIR}/${BIN}" "${GITDIR}/site/*" "${GITDIR}/cinder/*" "*/.*"

binzip: clean ${BIN}
	cd .. && zip ${GITDIR}-${VERSION}-$(shell uname -s).zip ${GITDIR}/${BIN}

.PHONY: www
