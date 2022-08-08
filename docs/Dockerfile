FROM golang:1.18.3 AS development
WORKDIR /ecommers
COPY go.mod go.sum ./
RUN go install github.com/cespare/reflex@latest
EXPOSE 8080
CMD reflex -g "*.go" go run cmd/main.go --start-service