FROM golang:1.20-alpine@sha256:87d0a3309b34e2ca732efd69fb899d3c420d3382370fd6e7e6d2cb5c930f27f9

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
