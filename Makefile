BUILD_DIR=bin
BUILD_NAME=servicetitan-to-dataset
BUILD_PREFIX=${BUILD_DIR}/${BUILD_NAME}

VERSION=0.1.0
LDFLAGS="-X servicetitan-to-dataset/cmd.version=$(VERSION)"

build:
	@mkdir -p ${BUILD_DIR}
	@GOOS=darwin GOARCH=amd64 go build -o ${BUILD_PREFIX}-darwin-amd64 -ldflags=${LDFLAGS}

test:
	@go test ./... -cover
