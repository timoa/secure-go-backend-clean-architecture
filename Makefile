build: bazel-build docker-build

bazel-build:
	bazelisk build //cmd:main

test: bazel-test

bazel-test:
	bazelisk test //...

coverage:
	go test -cover -coverprofile=coverage.out ./...

gazelle:
	bazelisk run //:gazelle

docker-build:
	docker-compose build

# Onboarding
onboarding: install-deps-macos setup

install-deps-macos:
	brew install pre-commit hadolint checkov gosec bazelisk

setup:
	pre-commit install

code-check:
	pre-commit run --all-files
