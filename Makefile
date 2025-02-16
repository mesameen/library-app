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

.PHONY: dockerbuild
dockerbuild:
	docker build -t mesameen/library-app .

.PHONY: dockerpush
dockerpush:
	docker push mesameen/library-app

.PHONY: dockerdeploy
dockerdeploy:
	docker compose up -d
