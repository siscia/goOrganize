package goorganizer

import (
	"fmt"
	"net/http"
	"time"
	"appengine"
	"github.com/gorilla/mux"
//	"appengine/datastore"
)


func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/p", handler1)
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	fmt.Fprint(w, time.Now(), "\n")
	NewThread(c, "simone@mweb.biz", "prova", "la1la")
	//if err != nil{fmt.Fprint(w, "???")}
	//NewPost(c, thread.Id, "simone@mweb.biz", "ole")
	//NewPost(c, thread.Id, "simone@mweb.biz", "sim")
	//GetUser(c, "simone@mweb.biz")
	//fmt.Fprint(w, thread.Participant)
}


func handler1(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	fmt.Fprint(w, time.Now(), "\n")
	NewThread(c, "leonardo@mweb.biz", "noidea", "???")
	//if err != nil{fmt.Fprint(w, "???")}
	//NewPost(c, thread.Id, "simone@mweb.biz", "ole")
	//NewPost(c, thread.Id, "simone@mweb.biz", "sim")
	//GetUser(c, "simone@mweb.biz")
	//fmt.Fprint(w, thread.Participant)
}