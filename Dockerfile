FROM golang:latest

RUN mkdir /opt/app

COPY . /opt/app
WORKDIR /opt/app

RUN go mod download

ENTRYPOINT ["/bin/bash", "-c"]
CMD ["go run ./main.go"]