package goorganizer

import (
	"time"
)

type User struct {
	Id string //id of own key
	Email string //id and email are the same...
	Nickname string
	Verified bool
	Password string
	FollowedThread []string //ids of thread.key
}

type RenderingThread struct{
	Thread
	EmailAuthor string
	Posts []Post
}

type Thread struct {
	Id string //id of own key
	Author string //id of key
	Title string
	Text string
	Time time.Time
	Posts []string //ids of Post.key
	Participant []string //id of User.key
}

type Post struct {
	Id string
	Author string //id of key
	Text string
	Time time.Time
}