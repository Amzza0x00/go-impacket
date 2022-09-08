BUILD_ENV := CGO_ENABLED=0
LDFLAGS=-v -a -ldflags '-s -w' -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"

PSEXEC := psexec

.PHONY: all setup build-linux build-osx build-windows

all: setup build-linux build-osx build-windows

setup:
	mkdir -p build/linux
	mkdir -p build/osx
	mkdir -p build/windows

build-linux:
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/linux/${PSEXEC}-linux-amd64 cmd/psexec.go;
	${BUILD_ENV} GOARCH=386 GOOS=linux go build ${LDFLAGS} -o build/linux/${PSEXEC}-linux-x86 cmd/psexec.go;

build-osx:
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o build/osx/${PSEXEC}-darwin-amd64 cmd/psexec.go;

build-windows:
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o build/windows/${PSEXEC}-windows-amd64.exe cmd/psexec.go;
	${BUILD_ENV} GOARCH=386 GOOS=windows go build ${LDFLAGS} -o build/windows/${PSEXEC}-windows-x86 cmd/psexec.go;
