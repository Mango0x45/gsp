.POSIX:

PREFIX = /usr/local
DPREFIX = ${DESTDIR}${PREFIX}

all: gsp
gsp:
	go build

install:
	mkdir -p ${DPREFIX}/bin \
	         ${DPREFIX}/share/man/man1 \
	         ${DPREFIX}/share/man/man5
	cp ${target}   ${DPREFIX}/bin
	cp ${target}.1 ${DPREFIX}/share/man/man1
	cp ${target}.5 ${DPREFIX}/share/man/man5

dist:
	mkdir -p dist
	for os in darwin linux windows; do \
		for arch in amd64 arm64; do \
			GOARCH=$$arch GOOS=$$os go build -o dist/gsp-$$os-$$arch; \
		done; \
	done

test:
	go test ./...

patch:
	patch vendor/github.com/tdewolff/parse/v2/position.go patches/newlines.diff

clean:
	rm -rf dist
