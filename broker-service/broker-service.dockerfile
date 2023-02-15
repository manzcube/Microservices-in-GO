# # base go image COMMMENTED OUT BECAUSE OF MAKEFILE, THAT CREATED EXE ALREADY
# FROM golang:1.18-alpine as builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp
# # build a tiny docekr image

FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

CMD [ "/app/brokerApp" ]
