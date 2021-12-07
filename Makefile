.PHONY: clean test

GOROOT=${HOME}/goroot
GOPATH=${HOME}/go2
PKG=github.com/dairaga/gs

test:
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/gofmt -w .
	- env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover -gcflags -G=3 .
	- env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover -gcflags -G=3 ${PKG}/cbf
	- env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover -gcflags -G=3 ${PKG}/either
	- env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover -gcflags -G=3 ${PKG}/funcs
	- env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover -gcflags -G=3 ${PKG}/option
	- env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover -gcflags -G=3 ${PKG}/slices
	- env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go test -v -cover -gcflags -G=3 ${PKG}/try
	
tidy:
	env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go mod tidy

clean:
	- env GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOROOT}/bin/go clean -cache -testcache -x