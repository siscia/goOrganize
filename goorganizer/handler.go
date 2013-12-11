package goorganizer

import (
	"fmt"
	"net/http"
	"html/template"
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


// all the error here should redirect to some other page explaining what went wrong and how to recover.

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
	r.ParseForm()
	p := r.PostForm
	email, text := p.Get("email"), p.Get("text")
	c.Infof("\n\n\n\nall: %v, email: %v, text: %v", p, email, text)
	_, err := NewPost(c, threadId, email, text)
	if err != nil {
		fmt.Fprint(w, "Get some problem, try again ", err)}
	fmt.Fprint(w, threadId)
}

func ServeReq(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	thread, err := GetThread(c, mux.Vars(r)["TH_ID"])
	if err != nil{
		c.Infof("get some problem, %v", err)}
	t := RenderingThread{Thread: thread, EmailAuthor: thread.Author}
	for _, post := range thread.Posts{
		u, _ := GetPost(c, post)
		t.Posts = append(t.Posts, u)}
	temp, err := template.ParseFiles("/home/simo/goOrganize/goorganizer/templates/Kreative10/thread2.html")
	if err != nil{
		c.Infof("Error: %v", err)
		panic("Parsing template panic")}
	c.Infof("thread in handler: %v", t)
	temp.Execute(w, t)
}

func DeletePostReq(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Delete Post Request")
}


func AddPeopleReq(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	thread, err := GetThread(c, mux.Vars(r)["TH_ID"])
	if err != nil{
		c.Infof("Error retriving thread x adding people, ", err)}
	r.ParseForm()
	user, err := GetUser(c, r.PostForm.Get("email"))
	if err != nil{
		c.Infof("Error retriving user x adding people, ", err)
	}
	thread, user = AddParticipant(c, thread, user)
	fmt.Fprint(w, "%v added at the %v, now s/he can post", user, thread)
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