package main

import (
	"context"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/haileyyoon23/christmas-wish-box/web"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	return r
}

func main() {
	log.Print("Christmas Wish Box Activation")
	log.Print("version : ", version)

	http.Handle("/", sentryHandler.Handle(buildHandler()))

	addr := "127.0.0.1:8000"//fmt.Sprintf("%s:%s", env.GetEnv().Domain, env.GetEnv().AppPort)
	println(addr)


	s := &http.Server{
		Addr: addr,//fmt.Sprintf("%s:%s", env.GetEnv().Domain, env.GetEnv().AppPort),
	}

	shutdownCompletion := make(chan struct{})

	go shutdownOnSignal(s, shutdownCompletion)

	log.Print("Listening Start")
	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal("Listen Failed")
	} else {
		<-shutdownCompletion
	}
}

func shutdownOnSignal(s *http.Server, completion chan<- struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Print("main: shutting down")
	err := s.Shutdown(ctx)
	if err != nil {
		log.Print("main: Shutdown() failed: ", err)
	}

	close(completion)
}