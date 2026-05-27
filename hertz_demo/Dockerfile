FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o /app/main .

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

ENV TZ=Asia/Shanghai

EXPOSE 8888

CMD ["./main"]

