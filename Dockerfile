FROM golang:1.14 as builder

RUN apt-get update
RUN apt-get install netcat -y

WORKDIR /app
COPY . /app

RUN go get ./...
RUN CGO_ENABLED=1 GOOS=linux go build -o rooms .
RUN chmod +x run.sh

ENTRYPOINT ["./run.sh"]