FROM golang:latest

RUN mkdir /opt/app

COPY . /opt/app
WORKDIR /opt/app

RUN go mod download

RUN go build -o christmas-wish-box .

#ENTRYPOINT ["/bin/bash", "-c"]
CMD ["/christmas-wish-box"]