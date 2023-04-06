FROM docker.io/library/golang:1.18-alpine as buildStage

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 go build -v -o list-pods .

FROM registry.obmondo.com/obmondo/dockerfiles/ubuntu:22.04

RUN apt-get update && apt install -y ca-certificates

WORKDIR /app

COPY --from=buildStage /app/list-pods /app/

CMD [ "/app/list-pods" ]
