# multi stage로 빌드하면 도커 이미지 사이즈가 줄어듬

#Build stage
FROM golang:1.18.5-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]

