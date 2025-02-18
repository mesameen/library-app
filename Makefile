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

.PHONY: dockerbuild
dockerbuild:
	docker build -t mesameen/library-app .

.PHONY: dockerpush
dockerpush:
	docker push mesameen/library-app

.PHONY: helminstall
helminstall:
	kubectl config use-context minikube
	helm3 upgrade --install library-app \
	--recreate-pods -f k8s/values.yaml \
	--set image.tag=latest \
	--namespace library-operations ./k8s

.PHONY: helmdeploy
helmdeploy:	dockerbuild dockerpush helminstall
