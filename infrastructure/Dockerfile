## Build
FROM golang:1.20-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o cli cmd/api/main.go


## Deploy
FROM alpine:3.18

# Set Timezone
RUN apk add tzdata
ENV TZ=Asia/Jakarta

WORKDIR /app

COPY /config ./config
COPY --from=builder /app/cli ./

EXPOSE 3000

ENTRYPOINT ["./cli"]