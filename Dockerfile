FROM golang:1.17-alpine as build-env

WORKDIR /

COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go build -o mutating-webhook

COPY --from=build-env /mutating-webhook /usr/local/bin/mutating-webhook

ENTRYPOINT ["mutating-webhook"]