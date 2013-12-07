package goorganizer

import (
	"fmt"
	"net/http"
	"time"
	"appengine"
)


func init() {
    http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	fmt.Fprint(w, time.Now())
	NewThread(c, "simone@mweb.biz", "prova", "lala")
	_, err := GetUser(c, "simone@mweb.biz")
	fmt.Fprint(w, "\nHello, world!\n")
	if err != nil{
		fmt.Fprint(w, err)}
}

