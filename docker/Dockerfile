FROM golang:1.26.2-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY hertz_demo/go.mod hertz_demo/go.sum ./hertz_demo/

RUN go work sync

COPY pkg/ ./pkg/
COPY hertz_demo/ ./hertz_demo/

WORKDIR /build/hertz_demo
RUN go build -ldflags="-w -s" -o /app/main .

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/main .
COPY --from=builder /build/hertz_demo/configs ./configs

ENV TZ=Asia/Shanghai

EXPOSE 8080

CMD ["./main"]

