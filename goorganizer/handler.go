package goorganizer

import (
	"fmt"
	"net/http"
	"time"
	"appengine"
	"github.com/gorilla/mux"
	"github.com/hoisie/mustache"
)


func init() {
	r := mux.NewRouter()
	r.HandleFunc("/handle-debug", DebugReq)
	r.HandleFunc("/new", NewThreadReq)
	r.HandleFunc("/new-post/{TH_ID}", NewPostReq)
	r.HandleFunc("/serve/{TH_ID}", ServeReq)
	r.HandleFunc("/delete-post/{TH_ID}", DeletePostReq)
	r.HandleFunc("/add-people/{TH_ID}", AddPeopleReq)
	r.HandleFunc("/modify-post/{TH_ID}", ModifyPostReq)
	r.HandleFunc("/get-json/{TH_ID}", GetJsonReq)
	r.HandleFunc("/", Main)
	http.Handle("/", r)
}

func DebugReq(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	r.ParseForm()
	req := r.PostForm.Get("request-type")
	fmt.Fprint(w, r.PostForm)
	switch req {
	case "new-thread": 
		NewThread( c,
			r.PostForm.Get("email"),
			r.PostForm.Get("title"),
			r.PostForm.Get("text")) 
	case "new-post":
		NewPost( c,
			string(r.PostForm.Get("TH_ID")),
			r.PostForm.Get("email"),
			r.PostForm.Get("text"))}
}

func ShowThread(w http.ResponseWriter, r *http.Request){
	
}

func NewThreadReq(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	p := r.Form
	email, title, text := p.Get("email"), p.Get("title"), p.Get("text")
	_, err := NewThread(c, email, title, text)
	if err != nil{
		fmt.Fprint(w, "Get some problem, try again", err)}
	fmt.Fprint(w, email, title, text)
}

func NewPostReq(w http.ResponseWriter, r *http.Request){
	threadId := mux.Vars(r)["TH_ID"]
	c := appengine.NewContext(r)
	p := r.Form
	email, text := p.Get("email"), p.Get("text")
	_, err := NewPost(c, threadId, email, text)
	if err != nil {
		fmt.Fprint(w, "Get some problem, try again", err)}
	fmt.Fprint(w, threadId)
}

func ServeReq(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Serve Thread Request")
}

func DeletePostReq(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Delete Post Request")
}


func AddPeopleReq(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Add People Request")
}


func ModifyPostReq(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Modify Post Request")
}


func GetJsonReq(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Get Json Request")
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	fmt.Fprint(w, time.Now(), "\n")
	NewThread(c, "simone@mweb.biz", "prova", "la1la")
}


func Main(w http.ResponseWriter, r *http.Request) {
	form := mustache.RenderFile("/home/simo/goOrganize/goorganizer/templates/main.mustache.html", nil)
	fmt.Fprint(w, form)
}