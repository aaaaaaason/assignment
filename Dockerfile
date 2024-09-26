FROM golang:1.23.1-alpine3.20 AS build

ENV GOOS=linux

WORKDIR build

COPY ./ .

RUN go get 

RUN CGO_ENABLED=0 go build -o app

FROM golang:1.23.1-alpine3.20

COPY --from=build go/build/app bin/app

ENTRYPOINT ["app"]
