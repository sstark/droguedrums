BIN=		droguedrums
GITDIR=		${BIN}
VERSION=	1.0
BUILD_TIME=	$(shell date +%FT%T%z)
LDFLAGS=	-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"
MIDILIB=	portmidi
SRCFILES=	src/main.go \
			src/matrix.go \
			src/config.go \
			src/midi.go \
			src/genlanes.go \
			src/velocity.go \
			src/midi-${MIDILIB}.go
PREFIX=		/usr/local

${BIN}: ${SRCFILES} Makefile
	@echo building for ${MIDILIB}
	go build ${LDFLAGS} -o ${BIN} ${SRCFILES}

test:
	cd src && go test

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
