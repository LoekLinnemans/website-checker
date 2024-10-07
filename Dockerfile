# Build stage
FROM golang:1.23.1-alpine AS build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/main .

COPY --from=build /app/config.txt .

COPY --from=build /app/result.txt .

COPY --from=build /app/log.txt .

CMD ["./main"]
