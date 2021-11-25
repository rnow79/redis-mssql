# syntax=docker/dockerfile:1
FROM golang:latest
WORKDIR /app
COPY go.mod ./
COPY photo.go ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /photo
RUN rm -r /app
ENTRYPOINT [ "/photo" ]

