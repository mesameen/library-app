MAINPATH = cmd/main.go

.PHONY: run
run:
	go run $(MAINPATH)

.PHONY: test
test:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
