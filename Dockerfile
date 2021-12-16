FROM golang:latest

WORKDIR /opt/app

COPY . .

RUN go mod download

RUN go get \
    && go get github.com/getsentry/sentry-go \
    && github.com/gorilla/mux \
    && github.com/labstack/gommon \
    && github.com/mattn/go-sqlite3 

RUN go build -o main .

#ENTRYPOINT ["/bin/bash", "-c"]
CMD ["./main"]
