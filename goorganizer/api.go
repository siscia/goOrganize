package goorganizer

import (
	"appengine/datastore"
	"time"
	"appengine"
	"errors"
	"encoding/json"
	"hash/fnv"
	"fmt"
)

func GenerateHash(text string, title string, moment time.Time) string{
	h := fnv.New64a()
	t := fmt.Sprintf("%x%x%x", text, title, moment)
	h.Write([]byte(t))
	return fmt.Sprintf("%x", int64(h.Sum64()))
}


func GetUser(c appengine.Context, email string) (User, error){
	var user User 
	key := datastore.NewKey(c, "Users", email, 0, nil) 
	if datastore.Get(c, key, user) == datastore.ErrNoSuchEntity{
		//panic("Not found email")
		return CreateUser(email, c)}
	c.Infof("user.Id: %v", user.Id)
	return user, nil
}

func CreateUser(email string, c appengine.Context) (User, error) {
	key := datastore.NewKey(c, "Users", email, 0, nil)
	id := key.StringID()
	c.Infof("user create: %v", id)
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
	hash := GenerateHash(text, title, time.Now())
	key := datastore.NewKey(c, "Thread", hash, 0, nil)
	id := key.StringID()
	c.Infof("id of thread: %v", id)
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

func GetPost(c appengine.Context, id string) (Post, error) {
	key := datastore.NewKey(c, "Posts", id, 0, nil)
	var post Post
	if datastore.Get(c, key, post) == datastore.ErrNoSuchEntity{
		return Post{}, datastore.ErrNoSuchEntity}
	return post, nil
}

func NewPost(c appengine.Context, threadId string, email string, text string) (Thread, error){
	thread, err := GetThread(c, threadId)
	if err != nil {
		return Thread{}, err}
	user, err := GetUser(c, email)
	c.Infof("thread.Id: %v user.Id: %v", thread.Id, user.Id)
	for _, participant := range thread.Participant {
		if participant == user.Id {
			id := GenerateHash(email, text, time.Now())
			key := datastore.NewKey(c, "Posts", id, 0, nil)
			post := Post{Id: key.StringID(), Author: user.Id, Text: text, Time: time.Now()}
			_, err := datastore.Put(c, key, post)
			if err != nil{
				panic("Error writing Post")}
			return AddPost(c, thread, post)}}
	panic("Non AUTH")
	//return Thread{}, errors.New("Not Authenticate User")
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
	thread.Posts =  append(thread.Posts, post.Id)
	return UpdateThread(c, thread)
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
		key := datastore.NewKey(c, "Posts", postID, 0, nil)
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

func UpdateUser(c appengine.Context, user User) (User, error){
	key := datastore.NewKey(c, "Users", user.Id, 0, nil)
	_, err := datastore.Put(c, key, user)
	if err != nil {
		return User{}, err}
	return user, nil
}

func UpdateThread(c appengine.Context, thread Thread) (Thread, error){
	key := datastore.NewKey(c, "Threads", thread.Id, 0, nil)
	_, err := datastore.Put(c, key, thread)
	if err != nil{
		return Thread{}, err}
	return thread, nil
}

func UpdatePost(c appengine.Context, post Post) (Post, error){
	key := datastore.NewKey(c, "Posts", post.Id, 0, nil)
	_, err := datastore.Put(c, key, post)
	if err != nil {
		return Post{}, err}
	return post, nil}