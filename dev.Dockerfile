FROM golang:1.22.1-bullseye

RUN apt-get update \
    && apt-get install -y ca-certificates curl gnupg \
    && mkdir -p /etc/apt/keyrings

WORKDIR /app

RUN go install go.uber.org/mock/mockgen@latest && \
    go install github.com/air-verse/air@latest


RUN git config --global --add safe.directory /app

CMD ["air", "-c", ".air.toml"]
