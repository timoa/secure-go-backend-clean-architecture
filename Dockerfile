FROM golang:1.20-alpine@sha256:1e2917143ce7e7bf8d1add2ac5c5fa3d358b2b5ddaae2bd6f54169ce68530ef0

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
