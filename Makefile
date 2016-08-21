BIN=		droguedrums
VERSION=	0.9
BUILD_TIME=	$(shell date +%FT%T%z)
LDFLAGS=	-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"
MIDILIB=	portmidi
SRCFILES=	src/main.go \
			src/matrix.go \
			src/config.go \
			src/midi.go \
			src/midi-${MIDILIB}.go

${BIN}: ${SRCFILES} Makefile
	@echo building for ${MIDILIB}
	go build ${LDFLAGS} -o ${BIN} ${SRCFILES}

test:
	cd ${SRC} && go test

clean:
	rm ${BIN}
