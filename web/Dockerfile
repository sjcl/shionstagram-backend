FROM golang:1.18 AS build

ENV CGO_ENABLED=0
WORKDIR /root/build

COPY go.mod .
COPY go.sum .

RUN set -x && go mod download

COPY . .

RUN set -x && go build

FROM alpine:3.15

RUN apk add --no-cache openssl

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY --from=build \
        /root/build/web \
        /usr/bin/shionstagram-backend

ENTRYPOINT ["dockerize", "-wait", "tcp://db:3306", "-timeout", "60s"]
CMD ["/usr/bin/shionstagram-backend"]