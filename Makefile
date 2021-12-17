.PHONY: clean test

GOROOT=${HOME}/goroot
GOPATH=${HOME}/go2
PKG=github.com/dairaga/gs

test:
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/gofmt -w .
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover .
	- env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover ${PKG}/cbf
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover ${PKG}/either
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover ${PKG}/funcs
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover ${PKG}/option
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover ${PKG}/slices
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover ${PKG}/maps
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover ${PKG}/try
	
tidy:
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go mod tidy

clean:
	- env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go clean -cache -testcache -x