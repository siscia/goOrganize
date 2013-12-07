package goorganizer

import (
	"appengine/datastore"
	"time"
	"appengine"
	"errors"
)


func GetUser(c appengine.Context, email string) User{
	var user User 
	key := datastore.NewKey(c, "Users", email, 0, nil) //it is going to create different keys if in a different Context (c) ??
	if datastore.Get(c, key, user) == datastore.ErrNoSuchEntity{
		return CreateUser(email, c)}
	return user
}

func CreateUser(email string, c appengine.Context) User{
	key := datastore.NewKey(c, "Users", email, 0, nil)
	user := User{Id: *key, Email: email, Verified: false}
	_, err := datastore.Put(c, key, user)
	if err != nil {
		panic("Error in Write User")}
	return user
}


func NewThread(c appengine.Context, email string, title string, text string) Thread {
	author := GetUser(c, email)
	key := datastore.NewIncompleteKey(c, "Thread", nil)
	thread := Thread{Id: *key, Author: author, Title: title, Text: text, Time: time.Now(), Participant: []User{author}}
	_, err := datastore.Put(c, key, thread)
	if err != nil {
		panic("Error in Write Thread")}
	return thread
}

func GetThread(c appengine.Context, id string) (Thread, error) {
	key := datastore.NewKey(c, "Thread", id, 0, nil)
	var thread Thread
	if datastore.Get(c, key, thread) == datastore.ErrNoSuchEntity{
		return Thread{}, datastore.ErrNoSuchEntity}
	return thread, nil
}

func NewPost(c appengine.Context, threadId string, email string, text string) (Thread, error){
	thread, err := GetThread(c, threadId)
	if err != nil {
		return Thread{}, err}
	user := GetUser(c, email)
	for _, participant := range thread.Participant {
		if participant == user {
			post := Post{Author: user, Text: text, Time: time.Now()}
			return AddPost(c, thread, post), nil}}
	return Thread{}, errors.New("Not Authenticate User")
}

func AddPost(c appengine.Context, thread Thread, post Post) Thread{
	new_posts := append(thread.Posts, post)
	thread.Posts = new_posts
	_, err := datastore.Put(c, &thread.Id, thread)
	if err != nil {
		panic("Error updating a Post")}
	return thread
}

func UpdatePost(c appengine.Context, thread Thread) (Thread, error){
	_, err := datastore.Put(c, &thread.Id, thread)
	if err != nil{
		return Thread{}, err}
	return thread, nil
}

func DeletePost(c appengine.Context, thread Thread, index int) (Thread, error){
	if len(thread.Posts) > index {
		thread.Posts = append(thread.Posts[:index], thread.Posts[index+1:]...)
		return UpdatePost(c, thread)
	}
	return Thread{}, errors.New("Out of Index")
}