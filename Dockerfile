FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary

## envirorment variables

ENV DB_HOST=88.222.214.98
ENV DB_NAME=user
ENV DB_PORT=5432
ENV DB_PASS=ameliasrg
ENV DB_USER=postgres
ENV FORM=http://localhost:8080
ENV SECRET=aksdSasiaSIOpwk049323


ENTRYPOINT ["/app/binary"]