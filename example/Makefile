darwin:
	CC=o64-clang CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build --installsuffix cgo -ldflags="-s -w" systray_example.go

windows:
	GOOS=windows GOARCH=amd64 go build systray_example.go

all: darwin windows
