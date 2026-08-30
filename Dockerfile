# -- Build stage --------------------------------------------------------------
FROM golang:1.26-bookworm@sha256:8e8aa801e8417ef0b5c42b504dd34db3db911bb73dba933bd8bde75ed815fdbb AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY api ./api
COPY bootstrap ./bootstrap
COPY cmd ./cmd
COPY domain ./domain
COPY internal ./internal
COPY mongo ./mongo
COPY repository ./repository
COPY usecase ./usecase

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd

# -- Final image -------------------------------------------------------------
FROM debian:bookworm-slim@sha256:e899040a73d36e2b36fa33216943539d9957cba8172b858097c2cabcdb20a3e2

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
  && rm -rf /var/lib/apt/lists/*

RUN useradd -r -u 10001 -g nogroup app-user

WORKDIR /app
COPY --from=builder /app/main /app/main
RUN chown -R app-user:nogroup /app

USER app-user

CMD ["/app/main"]
