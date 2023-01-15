build: go-build docker-build

docker-build:
	docker-compose build

go-test:
	go test -cover -coverprofile=coverage.out ./...

go-build: go-fmt go-get
	go build -o bin/main cmd/main.go

go-fmt:
	go fmt ./...

go-get:
	go get ./...

# Onboarding
onboarding: install-deps-macos setup

install-deps-macos:
	brew install pre-commit hadolint checkov gosec

setup:
	pre-commit install

code-check:
	pre-commit run --all-files
