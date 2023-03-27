FROM golang:1.20-alpine@sha256:576da1aa73f8ffa3dc6a7577b6032bd834aa84f2e1714d3e7e96b06b49f4e177

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
