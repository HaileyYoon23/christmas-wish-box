package web

import (
	"fmt"
	"github.com/haileyyoon23/christmas-wish-box/content"
	"github.com/haileyyoon23/christmas-wish-box/db"
	"net/http"
	"path"
	"runtime"
	"strconv"
)

type templateStruct map[string]interface{}

func HomePage(w http.ResponseWriter, r *http.Request) {
	gifts, err := db.GetGift(db.DB)
	if err != nil {
		panic(err)
	}

	errMsg := r.URL.Query().Get("errMsg")


	content.ExecuteTemplate(w, "home", templateStruct{
		"giftList": gifts,
		"errMsg": errMsg,
	})
}

func GiftAppendHandler(w http.ResponseWriter, r *http.Request) {
	gift := r.URL.Query().Get("gift")

	var path string

	err := db.AddGift(db.DB, gift)
	if err != nil {
		path = "?errMsg=" + err.Error()
	}

	http.Redirect(w, r, "/home" + path, http.StatusSeeOther)
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				renderInternalError(w, e)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func renderInternalError(w http.ResponseWriter, err interface{}) {
	_, file, line, _ := runtime.Caller(3)

	b := []byte(fmt.Sprintf("ERROR\n\n%s:%d\n\n%s\n", path.Base(file), line, err))

	h := w.Header()
	h.Set("Content-Type", "text/plain;charset=utf-8")
	h.Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(b)
}