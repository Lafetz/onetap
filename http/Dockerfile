FROM golang:1.22.3-alpine3.18 as builder


ENV APP_HOME=/go/src/http

WORKDIR "${APP_HOME}"


COPY ./go.mod ./go.sum ./ 


RUN go mod download
RUN go mod verify



COPY ./internal ./internal

COPY ./cmd ./cmd




RUN go build -o ./bin/http ./cmd

FROM alpine:latest

ENV APP_HOME=/go/src/http
RUN mkdir -p "${APP_HOME}"

WORKDIR "${APP_HOME}"

COPY --from=builder "${APP_HOME}"/bin/http $APP_HOME

EXPOSE 8080

CMD ["./http"]