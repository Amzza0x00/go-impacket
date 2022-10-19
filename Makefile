BUILD_ENV := CGO_ENABLED=0
LDFLAGS=-v -a -ldflags '-s -w' -gcflags="all=-trimpath=${PWD};${GOPATH};${GOROOT}" -asmflags="all=-trimpath=${PWD};${GOPATH};${GOROOT}"

PSEXEC := psexec
OXIDFIND := oxidfind

.PHONY: all setup build-linux build-osx build-windows

all: setup build-linux build-osx build-windows

setup:
	mkdir -p build/linux
	mkdir -p build/osx
	mkdir -p build/windows

build-linux:
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/linux/${PSEXEC}-linux-amd64 cmd/psexec.go;
	${BUILD_ENV} GOARCH=386 GOOS=linux go build ${LDFLAGS} -o build/linux/${PSEXEC}-linux-x86 cmd/psexec.go;
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/linux/${OXIDFIND}-linux-amd64 cmd/oxidfind.go;
	${BUILD_ENV} GOARCH=386 GOOS=linux go build ${LDFLAGS} -o build/linux/${OXIDFIND}-linux-x86 cmd/oxidfind.go;

build-osx:
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o build/osx/${PSEXEC}-darwin-amd64 cmd/psexec.go;
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o build/osx/${OXIDFIND}-darwin-amd64 cmd/oxidfind.go;


build-windows:
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o build/windows/${PSEXEC}-windows-amd64.exe cmd/psexec.go;
	${BUILD_ENV} GOARCH=386 GOOS=windows go build ${LDFLAGS} -o build/windows/${PSEXEC}-windows-x86.exe cmd/psexec.go;
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o build/windows/${OXIDFIND}-windows-amd64.exe cmd/oxidfind.go;
	${BUILD_ENV} GOARCH=386 GOOS=windows go build ${LDFLAGS} -o build/windows/${OXIDFIND}-windows-x86.exe cmd/oxidfind.go;

