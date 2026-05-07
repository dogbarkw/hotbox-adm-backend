FROM hub.aiera.tech/base_image/golang:1.16.6-alpine as builder
MAINTAINER Razil "woshilijinghua@gmail.com"
WORKDIR /app
ADD . .
RUN go mod download github.com/go-playground/validator/v10  && go get github.com/gin-gonic/gin/binding@v1.8.1 && go build -o /app/cmd
FROM hub.aiera.tech/base_image/alpine:3.11.6
WORKDIR /app
COPY --from=builder /app/cmd /app/cmd
COPY --from=builder /app/config /app/config
EXPOSE 8888