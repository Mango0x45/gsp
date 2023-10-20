.POSIX:

PREFIX = /usr/local
DPREFIX = ${DESTDIR}${PREFIX}

target = gsp
sources = main.go \
          formatter/formatter.go \
          parser/errors.go \
          parser/parser.go \
          parser/reader.go

all: ${target}
gsp: ${sources}
	go build

install:
	mkdir -p ${DPREFIX}/bin \
	         ${DPREFIX}/share/man/man1 \
	         ${DPREFIX}/share/man/man5
	cp ${target}   ${DPREFIX}/bin
	cp ${target}.1 ${DPREFIX}/share/man/man1
	cp ${target}.5 ${DPREFIX}/share/man/man5

test:
	go test . ./parser ./formatter
