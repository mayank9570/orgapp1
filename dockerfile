# syntax=docker/dockerfile:1



FROM golang:1.16-alpine



WORKDIR /organization-api
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o ./organization-api

CMD [ "./organization-api" ]