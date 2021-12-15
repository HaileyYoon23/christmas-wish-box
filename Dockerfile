FROM golang:latest

WORKDIR /src/app

COPY . .

RUN go mod download

ENTRYPOINT ["/bin/bash", "-c"]
CMD ["go run /src/app/main.go"]