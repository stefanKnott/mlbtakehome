.PHONY: fmt
fmt:
	go fmt ./...


.PHONY: test
test:
	go test --race ./...

.PHONY: run
run:
	go run mlbtakehome.go