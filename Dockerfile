FROM golang:1.25.4-alpine AS builder

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Kuala_Lumpur /etc/localtime && \
    echo "Asia/Kuala_Lumpur" > /etc/timezone

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o project .

FROM alpine:latest

RUN apk add --no-cache tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Kuala_Lumpur /etc/localtime && \
    echo "Asia/Kuala_Lumpur" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/project .

ENV TZ=Asia/Kuala_Lumpur

EXPOSE 8080
ENTRYPOINT ["./project"]