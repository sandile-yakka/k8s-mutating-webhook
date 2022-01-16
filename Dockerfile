FROM golang:1.17-alpine as build-env

WORKDIR /app

COPY go.mod /app
COPY go.sum /app

RUN go mod download

COPY . /app

RUN go build -o mutating-webhook

FROM golang:1.17-alpine

COPY --from=build-env /mutating-webhook /usr/local/bin/mutating-webhook

ENTRYPOINT ["mutating-webhook"]