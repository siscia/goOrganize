package goorganizer

import (
	"appengine/datastore"
	"time"
	"appengine"
	"errors"
	"encoding/json"
)


func GetUser(c appengine.Context, email string) (User, error){
	var user User 
	key := datastore.NewKey(c, "Users", email, 0, nil) 
	if datastore.Get(c, key, user) == datastore.ErrNoSuchEntity{
		return CreateUser(email, c)}
	return user, nil
}

func CreateUser(email string, c appengine.Context) (User, error) {
	key := datastore.NewKey(c, "Users", email, 0, nil)
	id := key.AppID()
	user := User{Id: id, Email: email, Verified: false}
	_, err := datastore.Put(c, key, &user)
	if err != nil {
		return User{}, err}
	return user, nil
}

func NewThread(c appengine.Context, email string, title string, text string) (Thread, error) {
	author, erro := GetUser(c, email)
	if erro != nil {
		return Thread{}, erro}
	key := datastore.NewIncompleteKey(c, "Thread", nil)
	id := key.AppID()
	thread := Thread{Id: id, Title: title, Text: text, Time: time.Now(), Participant: []string{author.Id}}
	_, err := datastore.Put(c, key, &thread)
	if err != nil {
		panic("Error in Write Thread")}
	return thread, nil
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
	user, err := GetUser(c, email)
	for _, participant := range thread.Participant {
		if participant == user.Id {
			key := datastore.NewIncompleteKey(c, "Post", nil)
			post := Post{Id: key.AppID(), Author: user.Id, Text: text, Time: time.Now()}
			_, err := datastore.Put(c, key, post)
			if err != nil{
				panic("Error writing Post")}
			return AddPost(c, thread, post)}}
	return Thread{}, errors.New("Not Authenticate User")
}

func AddParticipant(c appengine.Context, thread Thread, user User) (Thread, User) {
// all this should run in a transaction
	thread.Participant = append(thread.Participant, user.Id)
	t, err := UpdateThread(c, thread)
	if err != nil {
		panic("Error in adding a Participant")}
	user.FollowedThread = append(user.FollowedThread, thread.Id)
	u, err := UpdateUser(c, user)
	if err != nil{
		panic("Error in adding a Conversation to the followed")}
	return t, u
}

func AddPost(c appengine.Context, thread Thread, post Post) (Thread, error){
	new_posts := append(thread.Posts, post.Id)
	thread.Posts = new_posts
	return UpdateThread(c, thread)
}

func UpdateUser(c appengine.Context, user User) (User, error){
	key := datastore.NewKey(c, "User", user.Id, 0, nil)
	_, err := datastore.Put(c, key, user)
	if err != nil {
		return User{}, err}
	return user, nil
}

func UpdateThread(c appengine.Context, thread Thread) (Thread, error){
	key := datastore.NewKey(c, "Thread", thread.Id, 0, nil)
	_, err := datastore.Put(c, key, thread)
	if err != nil{
		return Thread{}, err}
	return thread, nil
}

func UpdatePost(c appengine.Context, post Post) (Post, error){
	key := datastore.NewKey(c, "Post", post.Id, 0, nil)
	_, err := datastore.Put(c, key, post)
	if err != nil {
		return Post{}, err}
	return post, nil
}

func DeletePost(c appengine.Context, thread Thread, index int) (Thread, error){
	if len(thread.Posts) > index {
		thread.Posts = append(thread.Posts[:index], thread.Posts[index+1:]...)
		return UpdateThread(c, thread)
	}
	return Thread{}, errors.New("Out of Index")
}

func ModifyPost(c appengine.Context, thread Thread, index int, text string) (Post, error){
	if len(thread.Posts) > index{
		postID := thread.Posts[index]
		key := datastore.NewKey(c, "Post", postID, 0, nil)
		var post Post
		if datastore.Get(c, key, post) == datastore.ErrNoSuchEntity{
			return Post{}, datastore.ErrNoSuchEntity}
		post.Text = text
		return UpdatePost(c, post)
	}
	return Post{}, errors.New("Out of Index")}

func JsonThread(thread Thread) ([]byte, error){
	j, err := json.Marshal(thread)
	if err != nil {
		return []byte{}, err}
	return j, nil}