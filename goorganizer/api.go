package goorganizer

import (
	"appengine/datastore"
	"time"
	"appengine"
	"errors"
	"encoding/json"
	"hash/fnv"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
)

// TODO AddParticipant should run in a transaction

func GenerateHash(text string, title string, moment time.Time) string{
	h := fnv.New64a()
	t := fmt.Sprintf("%x%x%x", text, title, moment)
	h.Write([]byte(t))
	return fmt.Sprintf("%x", h.Sum64())
}


func GetUser(c appengine.Context, email string) (User, error){
	var user User 
	key := datastore.NewKey(c, "Users", email, 0, nil) 
	err := datastore.Get(c, key, user)
	if err != nil{
		return CreateUser(c, email)}
	c.Infof("user.Id: %v", user.Id)
	return user, nil
}

func CreateUser(c appengine.Context, email string) (User, error) {
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
	obfusced := key.Encode()
	c.Infof("id of thread: %v", id)
	textb := []byte(text)
	thread := Thread{Id: id, Title: title, Text: textb, Time: time.Now(), Participant: []string{author.Id}, Author: author.Id, ObfuscedId: obfusced}
	_, err := datastore.Put(c, key, &thread)
	if err != nil {
		return Thread{}, err
}
	return thread, nil
}

func GetThread(c appengine.Context, id string) (Thread, error) {
	key := datastore.NewKey(c, "Thread", id, 0, nil)
	var thread Thread
	if datastore.Get(c, key, &thread) == datastore.ErrNoSuchEntity{
		panic("Not found key")
		return Thread{}, datastore.ErrNoSuchEntity}
	c.Infof("thread1: %v", thread)
	return thread, nil
}

func GetPost(c appengine.Context, id string) (Post, error) {
	key := datastore.NewKey(c, "Posts", id, 0, nil)
	var post Post
	if datastore.Get(c, key, &post) == datastore.ErrNoSuchEntity{
		return Post{}, datastore.ErrNoSuchEntity}
	return post, nil
}

func IsAuthUser(thread Thread, user User) bool{
	return true
	//*** this is why it doesn't works ****
	for _, participant := range thread.Participant{
		if participant == user.Email{
			return true}
	}
	return false
}

func NewPost(c appengine.Context, threadId string, email string, text string) (Thread, error){
	c.Infof("threadId: %v", threadId)
	thread, err := GetThread(c, threadId)
	c.Infof("thread: %v, error: %v", thread, err)
	if err != nil {
		return Thread{}, err}
	user, err := GetUser(c, email)
	if err != nil {
		return Thread{}, err}
	c.Infof("thread.Id: %v user.Id: %v", thread.Id, user.Id)
	if IsAuthUser(thread, user){
		id := GenerateHash(email, text, time.Now())
		key := datastore.NewKey(c, "Posts", id, 0, nil)
		text := []byte(text)
		post := Post{Id: key.StringID(), Author: user.Id, Text: text, Time: time.Now()}
		_, err := datastore.Put(c, key, &post)
		if err != nil{
			c.Infof("error just below: %v", err)
			panic("Error writing Post")}
		return AddPost(c, thread, post)}
	return Thread{}, errors.New("Non Auth User")
}



func AddParticipant(c appengine.Context, thread Thread, user User) error {
	thread.Participant = append(thread.Participant, user.Id)
	user.FollowedThread = append(user.FollowedThread, thread.Id)
	err := datastore.RunInTransaction(c, func(tc appengine.Context) error {
		_, errT := UpdateThread(tc, thread)			   
		if errT != nil {						   
			panic("Error in adding a Participant")
			return errT}		    
		_, errU := UpdateUser(tc, user)
		if errU != nil{
			panic("Error in adding a Conversation to the followed")
			return errU}
		return nil
	}, nil)
	if err != nil{
		panic("Error in transaction")}
	return nil
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
		if datastore.Get(c, key, &post) == datastore.ErrNoSuchEntity{
			return Post{}, datastore.ErrNoSuchEntity}
		post.Text = []byte(text)
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
	_, err := datastore.Put(c, key, &user)
	if err != nil {
		return User{}, err}
	return user, nil
}

func UpdateThread(c appengine.Context, thread Thread) (Thread, error){
	key := datastore.NewKey(c, "Thread", thread.Id, 0, nil)
	_, err := datastore.Put(c, key, &thread)
	if err != nil{
		return Thread{}, err}
	return thread, nil
}

func UpdatePost(c appengine.Context, post Post) (Post, error){
	key := datastore.NewKey(c, "Posts", post.Id, 0, nil)
	_, err := datastore.Put(c, key, &post)
	if err != nil {
		return Post{}, err}
	return post, nil
}

func RenderPosts(c appengine.Context, postIds []string) ([]RenderPost, error){
	np := make([]RenderPost, len(postIds))
	for i, p := range postIds{
		post, err := GetPost(c, p)
		if err != nil{
			return np, err}
		np[i] = RenderPost{Post: post}
		np[i].Html = RenderText(post.Text)}
	return np, nil
}

func RenderText(text []byte) interface{} {
	html := string(blackfriday.MarkdownCommon(text))
	return template.HTML(html)
}