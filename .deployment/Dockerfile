FROM golang:1.10.3
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/nmrshll/gphotos-uploader-cli

COPY Gopkg.lock Gopkg.toml ./
COPY vendor vendor
COPY main.go main.go
COPY config config
COPY datastore datastore
COPY filetypes filetypes
COPY gphotosapiclient gphotosapiclient
COPY upload upload

RUN go install ./...