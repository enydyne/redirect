FROM golang:latest as build
COPY . /app
WORKDIR /app
RUN go mod tidy
RUN go build -o main .
ENTRYPOINT ["/app/main"]
