FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go build -o qod .

FROM scratch
COPY --from=builder /app/qod /qod
EXPOSE 8080
ENTRYPOINT ["/qod"]
