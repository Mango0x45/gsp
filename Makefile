.POSIX:

PREFIX = /usr/local
DPREFIX = ${DESTDIR}${PREFIX}

all: gsp
gsp:
	go build

install:
	mkdir -p ${DPREFIX}/bin                                                     \
	         ${DPREFIX}/share/man/man1                                          \
	         ${DPREFIX}/share/man/man5                                          \
	         ${DPREFIX}/share/doc/gsp
	cp gsp ${DPREFIX}/bin
	cp man/gsp.1 ${DPREFIX}/share/man/man1
	sed 's#@DOCPATH@#${DPREFIX}/share/doc/gsp#' man/gsp.5 \
		>${DPREFIX}/share/man/man5/gsp.5
	cp example.gsp ${DPREFIX}/share/doc/gsp

dist:
	mkdir -p dist
	for os in darwin linux windows; do                                          \
		for arch in amd64 arm64; do                                             \
			GOARCH=$$arch GOOS=$$os go build -o dist/gsp-$$os-$$arch;           \
		done;                                                                   \
	done

test:
	go test ./...

patch:
	patch vendor/github.com/tdewolff/parse/v2/position.go patches/newlines.diff

clean:
	rm -rf dist

.PHONY: all clean dist gsp install patch test
