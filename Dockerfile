FROM golang:1.14-alpine as env

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0

FROM env as dev

RUN echo 'alias ll="ls -lah"' >> ~/.bashrc

FROM env as base

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go mod vendor

FROM base as builder

RUN go build -o api

FROM env as final

CMD ["/path/to/the/app", "root"]
