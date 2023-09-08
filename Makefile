PREFIX = /usr/local
DPREFIX = ${DESTDIR}/${PREFIX}

target = gsp

all: ${target}
gsp:
	go build

install:
	mkdir -p ${DPREFIX}/bin ${DPREFIX}/share/man/man1
	cp ${target}   ${DPREFIX}/bin
	cp ${target}.1 ${DPREFIX}/share/man/man1
