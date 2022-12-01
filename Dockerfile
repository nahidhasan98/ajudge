# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine

WORKDIR /app

COPY . .

# RUN go mod download

RUN go build -o ajudge .

# EXPOSE 8080

CMD [ "/app/ajudge" ]

