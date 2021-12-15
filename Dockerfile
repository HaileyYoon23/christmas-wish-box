FROM golang:latest

RUN mkdir /opt/app

COPY . /opt/app
WORKDIR /opt/app

RUN go mod download

RUN go build -o main .

#ENTRYPOINT ["/bin/bash", "-c"]
CMD ["/main"]