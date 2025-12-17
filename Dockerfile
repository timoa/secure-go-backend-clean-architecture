FROM alpine:3.21 AS builder

RUN apk add --no-cache \
    bash=5.2.37-r0 \
    ca-certificates=20250911-r0 \
    curl=8.14.1-r2 \
    git=2.47.3-r0

WORKDIR /app
COPY .bazelrc .bazelversion BUILD.bazel MODULE.bazel MODULE.bazel.lock go.mod go.sum ./
COPY api ./api
COPY bootstrap ./bootstrap
COPY cmd ./cmd
COPY domain ./domain
COPY internal ./internal
COPY mongo ./mongo
COPY repository ./repository
COPY usecase ./usecase

# Install bazelisk (respects .bazelversion)
RUN curl -L -o /usr/local/bin/bazelisk https://github.com/bazelbuild/bazelisk/releases/download/v1.27.0/bazelisk-linux-amd64 && \
    chmod +x /usr/local/bin/bazelisk

# Build the binary using Bazel
RUN bazelisk build //cmd:main

FROM alpine:3.21

RUN adduser -S app-user

WORKDIR /app
COPY --from=builder /app/bazel-bin/cmd/main /app/main
RUN chown -R app-user /app

USER app-user

CMD ["/app/main"]
