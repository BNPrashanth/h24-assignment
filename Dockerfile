FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add git
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-a -installsuffix cgo \
	-o h24-analyser-service ./cmd/main.go


FROM alpine
WORKDIR /app
COPY --from=builder /app/ /app/
EXPOSE 8080
CMD ["./h24-analyser-service"]
