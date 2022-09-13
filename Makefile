GOBIN 	:= go
PWD 	:= $(shell pwd)
BINDIR 	:= ${PWD}/bin
CFGDIR 	:= ${PWD}/conf
TARGET 	:= go-transfer

all:
	${GOBIN} mod tidy; \
	CGO_ENABLED=0 ${GOBIN} build -tags=${ENV} -o ${BINDIR}/${TARGET} .;

