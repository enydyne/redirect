FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch
WORKDIR /
COPY --from=build /app/main .
COPY --from=build /app/urls.txt .
EXPOSE 8080
ENTRYPOINT ["/main"]
