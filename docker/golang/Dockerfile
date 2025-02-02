#
# builder
#
FROM golang:1.23-bullseye AS builder

WORKDIR /workspace
COPY . .
RUN go mod tidy
RUN go build -o ./server .

#
# server app
#
FROM gcr.io/distroless/base-debian12:latest

WORKDIR /app
EXPOSE 8080
COPY --from=builder /workspace/server ./server
ENTRYPOINT ["/app/server"]
