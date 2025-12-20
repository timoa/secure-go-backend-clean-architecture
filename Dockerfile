# -- Build stage --------------------------------------------------------------
FROM alpine:3.22 AS builder

RUN apk add --no-cache \
    bash \
    ca-certificates \
    curl \
    git \
    unzip \
    build-base \
    gcompat

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
RUN TARGETARCH="${TARGETARCH:-amd64}" && \
    case "$TARGETARCH" in \
      "amd64")  BAZELISK_ARCH="amd64" ;; \
      "arm64")  BAZELISK_ARCH="arm64" ;; \
      *) echo "Unsupported TARGETARCH=$TARGETARCH"; exit 1 ;; \
    esac && \
    curl --fail --location --proto '=https' --proto-redir '=https' --tlsv1.2 \
      -o /usr/local/bin/bazelisk \
      "https://github.com/bazelbuild/bazelisk/releases/download/v1.27.0/bazelisk-linux-${BAZELISK_ARCH}" && \
    chmod +x /usr/local/bin/bazelisk

# Build the binary using Bazel
RUN bazelisk build //cmd:main

# -- Final image -------------------------------------------------------------
FROM alpine:3.22

RUN adduser -S app-user

WORKDIR /app
COPY --from=builder /app/bazel-bin/cmd/main /app/main
RUN chown -R app-user /app

USER app-user

CMD ["/app/main"]
