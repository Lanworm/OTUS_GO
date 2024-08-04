FROM golang:1.19

WORKDIR /opt/tests

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD go test -v ./tests/integration/...
