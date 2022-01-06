FROM golang:1.17.5-alpine3.15 as build
RUN apk add --update --no-cache git gcc musl-dev make


ENV GO111MODULE=on

WORKDIR /usr/github.com/BenjaminCallahan/my-bank-service

COPY . .

RUN go mod download

RUN go build -o ./.bin/web_app cmd/main.go


FROM alpine:3.15.0

RUN apk add ca-certificates

WORKDIR /usr/local/bin/app

COPY --from=build /usr/github.com/BenjaminCallahan/my-bank-service/.bin/web_app .
COPY schema ./schema/

EXPOSE 8080

CMD ["/usr/local/bin/app/web_app"]