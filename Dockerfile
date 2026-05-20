FROM golang:1.26

WORKDIR /app

RUN curl -sSf https://atlasgo.sh | sh
