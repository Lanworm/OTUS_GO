FROM golang:1.19

RUN  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0
COPY ../../build/app/migrations .
CMD migrate -database "postgres://postgres:postgres@calendar-postgres:5432/postgres?sslmode=disable" -path ./ up