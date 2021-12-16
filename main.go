package main

import (
	"fmt"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/haileyyoon23/christmas-wish-box/web"
	"github.com/labstack/gommon/log"
	"net/http"
)

var (
	version 	= "0.0.0-dev"
	sentryHandler = sentryhttp.New(sentryhttp.Options{Repanic: true})
)

func buildHandler() *mux.Router {
	r := mux.NewRouter()

	r.Use(web.ErrorHandler)

	r.HandleFunc("/home", web.HomePage).Methods("GET")
	r.HandleFunc("/index/add", web.GiftAppendHandler).Methods("GET")
	r.HandleFunc("/index/like", web.GiftLikeHandler).Methods("GET")
	return r
}

func main() {
	log.Print("Christmas Wish Box Activation")
	log.Print("version : ", version)

	http.Handle("/", sentryHandler.Handle(buildHandler()))

	fmt.Println("Server started at port 9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
