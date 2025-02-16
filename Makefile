MAINPATH = main.go

.PHONY: run
run:
	go run $(MAINPATH)

.PHONY: test
test:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

.PHONY: swag
swag:
	swag init

.PHONY: dockerdeploy
dockerdeploy:
	docker compose up -d
