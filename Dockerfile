FROM golang:1.17-alpine

WORKDIR /app
COPY . .

RUN go build
CMD ./netatmo-influxdb

