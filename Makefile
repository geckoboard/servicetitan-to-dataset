BUILD_DIR=bin
BUILD_NAME=servicetitan-to-dataset
BUILD_PREFIX=${BUILD_DIR}/${BUILD_NAME}

VERSION=0.1.0
LDFLAGS="-X servicetitan-to-dataset/cmd.version=$(VERSION)"

build:
	@mkdir -p ${BUILD_DIR}
	@GOOS=darwin  GOARCH=amd64  go build -o ${BUILD_PREFIX}-darwin-amd64  -ldflags=${LDFLAGS}
	@GOOS=darwin  GOARCH=arm64  go build -o ${BUILD_PREFIX}-darwin-arm64  -ldflags=${LDFLAGS}
	@GOOS=linux   GOARCH=386    go build -o ${BUILD_PREFIX}-linux-x86     -ldflags=${LDFLAGS}
	@GOOS=linux   GOARCH=amd64  go build -o ${BUILD_PREFIX}-linux-amd64   -ldflags=${LDFLAGS}
	@GOOS=windows GOARCH=386    go build -o ${BUILD_PREFIX}-windows-x86   -ldflags=${LDFLAGS}
	@GOOS=windows GOARCH=amd64  go build -o ${BUILD_PREFIX}-windows-amd64 -ldflags=${LDFLAGS}

test:
	@go test ./... -cover
