APP_NAME=aws-generate-secrets

.PHONY: build
build:
	CGO_ENABLED=0 go build \
		-o "build/$(APP_NAME)" \
		"./"

docker-build:
	docker build \
		-t "$(IMAGE)" \
		.