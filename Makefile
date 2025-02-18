NAMESPACE=library-operations
MAINPATH=main.go
GITTAG=$$(git tag --sort=committerdate | tail -1)

REGISTRY=docker.io/mesameen/library-app
IMG=$(REGISTRY):$(GITTAG)

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
	# docker build -t mesameen/library-app .
	docker build -t ${IMG}-dev .

.PHONY: dockerpush
dockerpush:
	# docker push mesameen/library-app
	docker push ${IMG}-dev

.PHONY: helminstall
helminstall:
	kubectl config use-context minikube
	helm3 upgrade --install library-app \
	--recreate-pods -f k8s/values.yaml \
	--set image.tag=${GITTAG}-dev \
	--namespace ${NAMESPACE} ./k8s

.PHONY: helmdeploy
helmdeploy:	dockerbuild dockerpush helminstall
