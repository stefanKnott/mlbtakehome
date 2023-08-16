.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test --race ./...

.PHONY: run
run:
	go mod download
	go run mlbtakehome.go

.PHONY: binaries
binaries:
	GOOS=darwin GOARCH=amd64 go build -o ./binaries/mlbtakehome_darwin_amd64 .
	GOOS=darwin GOARCH=arm64 go build -o ./binaries/mlbtakehome_darwin_arm64 .
	GOOS=linux GOARCH=amd64 go build -o ./binaries/mlbtakehome_linux_amd64 .
	GOOS=linux GOARCH=arm64 go build -o ./binaries/mlbtakehome_linux_arm64 .
	GOOS=windows GOARCH=amd64 go build -o ./binaries/mlbtakehome_windows_amd64 .
	GOOS=windows GOARCH=arm64 go build -o ./binaries/mlbtakehome_windows_arm64 .