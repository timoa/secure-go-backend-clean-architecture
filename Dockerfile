FROM debian:bookworm-slim AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    curl \
    git \
    unzip \
  && rm -rf /var/lib/apt/lists/*

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
RUN curl --fail --location --proto '=https' --proto-redir '=https' --tlsv1.2 -o /usr/local/bin/bazelisk https://github.com/bazelbuild/bazelisk/releases/download/v1.27.0/bazelisk-linux-amd64 && \
    chmod +x /usr/local/bin/bazelisk

# Build the binary using Bazel
RUN bazelisk build //cmd:main

FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

RUN adduser -S app-user

WORKDIR /app
COPY --from=builder /app/bazel-bin/cmd/main /app/main
RUN chown -R app-user /app

USER app-user

CMD ["/app/main"]
