FROM golang:1.19-alpine@sha256:2381c1e5f8350a901597d633b2e517775eeac7a6682be39225a93b22cfd0f8bb

# Create folder /app and non-privileged user as root
RUN mkdir /app && \
    adduser -S app-user

# Copy the project to the /app folder
COPY . /app

# Set the current folder as /app
WORKDIR /app

# Build the app
RUN go build -o main cmd/main.go && \
    chown -R app-user /app

# Use the non-privileged user for next actions
USER app-user

# Set the entrypoint
CMD ["/app/main"]
