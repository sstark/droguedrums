BIN=		droguedrums
GITDIR=		${BIN}
MIDILIB=	""
PREFIX=		/usr/local
OS=         $(shell uname)
VERSION=	$(shell sed '/version/!d;s/.*\"\(.*\)\"/\1/' version.go)

${BIN}: *.go Makefile
	go build -tags ${MIDILIB} -o ${BIN}

test:
	go test

test-all:
ifeq ($(OS),Darwin)
	go test -tags coremidi
endif
	go test -tags portmidi

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
