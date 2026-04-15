.POSIX:

PREFIX = /usr/local
DPREFIX = ${DESTDIR}${PREFIX}

all: build
build:
	printf 'package main\nconst installPrefix = "%s"\n' "${PREFIX}" \
		>config.gen.go
	go build

install:
	mkdir -p ${DPREFIX}/bin                                                     \
	         ${DPREFIX}/share/gsp                                               \
	         ${DPREFIX}/share/man/man1                                          \
	         ${DPREFIX}/share/man/man5
	cp gsp   ${DPREFIX}/bin
	cp gsp.1 ${DPREFIX}/share/man/man1
	cp gsp.5 ${DPREFIX}/share/man/man5
	cp -R macros ${DPREFIX}/share/gsp

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

.PHONY: all build clean dist install patch test
