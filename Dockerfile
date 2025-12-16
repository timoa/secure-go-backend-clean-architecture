FROM alpine:3.21 AS builder

RUN apk add --no-cache bash ca-certificates curl git

WORKDIR /app
COPY . .

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
