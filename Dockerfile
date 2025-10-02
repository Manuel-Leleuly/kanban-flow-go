# Builder
FROM golang:1.24.6-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

EXPOSE 8080


# DEV
FROM golang:1.24.6-alpine AS development

WORKDIR /app

RUN apk add --no-cache git

RUN go install github.com/air-verse/air@v1.62.0

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 3005

CMD ["air", "-c", ".air.toml"]

# PROD
FROM alpine:latest AS production

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080
ENV PORT=8080

CMD [ "./main" ]