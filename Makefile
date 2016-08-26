BIN=		droguedrums
VERSION=	0.9
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
	rm ${BIN}

checkfmt:
	@gofmt -d src/*.go

install: ${BIN}
	install ${BIN} ${PREFIX}/bin

www: mkdocs.yml
	mkdocs build --clean

.PHONY: www
