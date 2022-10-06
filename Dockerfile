FROM golang:alpine AS builder
WORKDIR /ecommers
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main main.go

# RUN go install github.com/cespare/reflex@latest
# EXPOSE 8080
# CMD reflex -g "*.go" go run cmd/main.go --start-service

FROM alpine
WORKDIR /ecommers
COPY --from=builder /ecommers/main .
COPY . .
COPY .env . 
COPY start.sh .
COPY wait-for.sh .


EXPOSE 8080
CMD [ "/ecommers/main" ]
ENTRYPOINT [ "/ecommers/start.sh" ]