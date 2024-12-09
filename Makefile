
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -ldflags="-X 'main.version=${GIT_REV}'" -ldflags="-X 'main.author=${CURRENT_DIR}'" ./