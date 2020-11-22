FROM golang

WORKDIR /go/src/gphotos-uploader-cli/
COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .
RUN \
  make test && \
  make build && \
  make clean

ENTRYPOINT [ "./gphotos-uploader-cli" ]
