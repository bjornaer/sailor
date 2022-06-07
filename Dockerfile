FROM golang:1.18-alpine as builder

ENV APP_NAME sailor
ENV CMD_PATH cmd/PortDomainService/

WORKDIR $GOPATH/src/$APP_NAME
COPY go.mod go.sum ./
RUN go mod download 
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

FROM alpine:3.16

ENV APP_NAME sailor
# COPY ./ports.json .

COPY --from=builder /$APP_NAME .

EXPOSE 8081

CMD ./$APP_NAME

