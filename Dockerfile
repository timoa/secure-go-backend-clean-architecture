FROM golang:1.20-alpine@sha256:913de96707b0460bcfdfe422796bb6e559fc300f6c53286777805a9a3010a5ea

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
